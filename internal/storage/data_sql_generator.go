package storage

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
	"strconv"
	"strings"
)

func getDataUpsertQuery(tableName string) string {

	numFields := dataFieldsCount[tableName]

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

func getDataDeleteQuery(tableName string) string {

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("DELETE FROM public.data_%s WHERE user_id = $1 AND name = $2", tableName))
	return sb.String()
}

func getMetaDataUpsertQuery(tableName string) string {
	return fmt.Sprintf(`INSERT INTO public.metadata_%s (data_id, name, content, md5)
		SELECT id, $3, pgp_sym_encrypt($4, $5), md5($4)
		FROM public.data_%s
		WHERE user_id = $1 AND name = $2`, tableName, tableName)

}

func getMetaDataDeleteQuery(tableName string) string {
	return fmt.Sprintf("DELETE FROM public.metadata_%s WHERE data_id IN (SELECT id FROM public.data_%s WHERE user_id = $1 AND name = $2 ) AND name = $3", tableName, tableName)
}

func getDataReadQuery(tableName string, request *proto.DataReadRequest, cfg *config.ServerConfig) (string, []interface{}) {
	params := make([]interface{}, 0)
	params = append(params, request.NameMask)
	params = append(params, cfg.Key)

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`SELECT data_%s.id, data_%s.name`, tableName, tableName))

	for i := 0; i < dataFieldsCount[tableName]; i++ {
		sb.WriteString(fmt.Sprintf(",\npgp_sym_decrypt(content_%d, $2) AS content_%d", i+1, i+1))
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
			sb.WriteString(fmt.Sprintf(`(pgp_sym_decrypt(metadata_%s."content", $2) = $`+strconv.Itoa(2*i+3)+` AND metadata_%s.name=$`+strconv.Itoa(2*i+4)+`)`, tableName, tableName))
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

	sb.WriteString(` name LIKE $1`)

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
