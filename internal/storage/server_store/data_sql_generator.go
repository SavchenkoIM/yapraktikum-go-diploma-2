package server_store

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
	"strconv"
	"strings"
)

func getDataUpsertQuery(objectType proto.DataType) string {

	numFields := objectTypes[objectType].FieldsCount
	tableName := objectTypes[objectType].TableName

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("INSERT INTO public.data_%s (user_id, name, ", tableName))
	for i := 0; i < int(numFields); i++ {
		sb.WriteString(fmt.Sprintf("content_%d, ", i+1))
	}
	sb.WriteString("md5) VALUES (\n\t$1,\n\t$2, \n")
	for i := 0; i < int(numFields); i++ {
		sb.WriteString(fmt.Sprintf("\tpgp_sym_encrypt($%d, $%d), \n", i+3, numFields+3))
	}
	sb.WriteString("\tmd5(")
	for i := 0; i < int(numFields); i++ {
		sb.WriteString(fmt.Sprintf(`$%d`, i+3))
		if i < int(numFields-1) {
			sb.WriteString(fmt.Sprintf(" || "))
		}
	}
	sb.WriteString("))\n")
	sb.WriteString("ON CONFLICT (user_id, name) DO UPDATE SET \n")
	for i := 0; i < int(numFields); i++ {
		sb.WriteString(fmt.Sprintf("\tcontent_%d = excluded.content_%d, \n", i+1, i+1))
	}
	sb.WriteString("\tmd5 = excluded.md5")
	return sb.String()
}

func getDataDeleteQuery(objectType proto.DataType) string {

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("DELETE FROM public.data_%s WHERE user_id = $1 AND name = $2", objectTypes[objectType].TableName))
	return sb.String()
}

func getMetaDataUpsertQuery(objectType proto.DataType) string {
	tableName := objectTypes[objectType].TableName
	return fmt.Sprintf(`INSERT INTO public.metadata_%s (data_id, name, content, md5)
		SELECT id, $3, pgp_sym_encrypt($4, $5), md5($4)
		FROM public.data_%s
		WHERE user_id = $1 AND name = $2`, tableName, tableName)

}

func getMetaDataDeleteQuery(objectType proto.DataType) string {
	tableName := objectTypes[objectType].TableName
	return fmt.Sprintf("DELETE FROM public.metadata_%s WHERE data_id IN (SELECT id FROM public.data_%s WHERE user_id = $1 AND name = $2 ) AND name = $3", tableName, tableName)
}

func getDataReadQuery(objectType proto.DataType, userID string, request *proto.DataReadRequest, cfg *config.ServerConfig) (string, []interface{}) {
	params := make([]interface{}, 0)
	params = append(params, userID)
	params = append(params, request.NameMask)
	params = append(params, cfg.Key)

	tableName := objectTypes[objectType].TableName
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`SELECT data_%s.id, data_%s.name`, tableName, tableName))

	for i := 0; i < objectTypes[objectType].FieldsCount; i++ {
		sb.WriteString(fmt.Sprintf(",\npgp_sym_decrypt(content_%d, $3) AS content_%d", i+1, i+1))
	}

	if len(request.Metadata) == 0 {
		sb.WriteString(fmt.Sprintf(` FROM data_%s WHERE`, tableName))
	} else {
		sb.WriteString(fmt.Sprintf(`
FROM data_%s WHERE EXISTS (
SELECT id FROM metadata_%s WHERE data_%s.id = metadata_%s.data_id `, tableName, tableName, tableName, tableName))

		if len(request.Metadata) > 0 {
			sb.WriteString(`AND (`)
		}
		for i, v := range request.Metadata {
			sb.WriteString(fmt.Sprintf(`(pgp_sym_decrypt(metadata_%s."content", $3) LIKE $`+strconv.Itoa(2*i+4)+` AND metadata_%s.name=$`+strconv.Itoa(2*i+5)+`)`, tableName, tableName))
			if i < len(request.Metadata)-1 {
				sb.WriteString(" OR")
			}
			sb.WriteString("\n")
			params = append(params, v.Value)
			params = append(params, v.Name)
		}
		if len(request.Metadata) > 0 {
			sb.WriteString(`)`)
		}
		sb.WriteString(`)  AND `)
	}

	sb.WriteString(` name LIKE $2 `)
	sb.WriteString(` AND user_id = $1 `)

	return sb.String(), params
}

func getMetadataReadQuery(tableName string, parentUUID pgtype.UUID, cfg *config.ServerConfig) (string, []interface{}) {
	params := make([]interface{}, 0)
	params = append(params, cfg.Key)
	params = append(params, parentUUID)

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`SELECT name, pgp_sym_decrypt(content, $1) FROM public.metadata_%s WHERE
		public.metadata_%s.data_id = $2`, tableName, tableName))

	return sb.String(), params
}

func getDataReadQueryFull(objectType proto.DataType, userID string, request *proto.DataReadRequest, cfg *config.ServerConfig) (string, []interface{}) {
	params := make([]interface{}, 0)
	params = append(params, userID)
	params = append(params, request.NameMask)
	params = append(params, cfg.Key)

	tableName := objectTypes[objectType].TableName

	contents := strings.Builder{}
	nm := map[bool]string{true: "\n", false: ""}
	for i := 0; i < objectTypes[objectType].FieldsCount; i++ {
		contents.WriteString(fmt.Sprintf("\tpgp_sym_decrypt(data.content_%d, $3) AS content_%d,%s", i+1, i+1, nm[i < objectTypes[objectType].FieldsCount-1]))
	}

	mdFilter := strings.Builder{}
	orm := map[bool]string{true: " OR", false: ""}
	if request.Metadata != nil && len(request.Metadata) > 0 {
		mdFilter.WriteString(fmt.Sprintf("  AND \n\tEXISTS ( SELECT id FROM metadata_%s AS metadata WHERE\n\t\tdata.id = metadata.data_id AND (\n", tableName))
		for i, v := range request.Metadata {
			v := v
			mdFilter.WriteString(fmt.Sprintf("\t\t(metadata.name = $%d AND\n", i*2+4))
			mdFilter.WriteString(fmt.Sprintf("\t\t\tpgp_sym_decrypt(metadata.content, $3) LIKE $%d)%s\n", i*2+5, orm[i < len(request.Metadata)-1]))
			params = append(params, v.Name)
			params = append(params, v.Value)
		}
		mdFilter.WriteString("))")
	}

	queryTmpl := fmt.Sprintf(
		`SELECT 
	data.name AS data_name,
%s
	metadata.name AS m_name,
	pgp_sym_decrypt(metadata.content, $3) AS m_content
FROM 
	(SELECT * FROM data_%s WHERE user_id = $1) AS data
	LEFT JOIN metadata_%s AS metadata ON 
		data.id = metadata.data_id 
WHERE
	data.name LIKE $2 %s
`, contents.String(), tableName, tableName, mdFilter.String())

	return queryTmpl, params
}
