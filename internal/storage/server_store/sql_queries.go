package server_store

import (
	"fmt"
	"strings"
)

var dataFieldsCount = map[string]int{
	"text_note":   1,
	"credentials": 2,
	"credit_card": 3,
	"blob":        1,
}

var queryCreateExtensionUUID = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`

var queryCreateExtensionPGCrypto = `CREATE EXTENSION IF NOT EXISTS "pgcrypto"`

var queryCreateUsers string = `CREATE TABLE IF NOT EXISTS public.users
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    login text NOT NULL,
    password text NOT NULL,
    salt text NOT NULL,
    filestore_access_key bytea DEFAULT NULL,
    PRIMARY KEY (id),
    CONSTRAINT uk_login UNIQUE (login)
)
WITH (
    OIDS = FALSE
);`

var queryCreateData string = `CREATE TABLE IF NOT EXISTS public.data_{table}
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name text NOT NULL,
	{content}
    md5 char(32),
    PRIMARY KEY (id),
    CONSTRAINT uk_data_{table}_name UNIQUE (user_id, name)
)
WITH ( 
    OIDS = FALSE
);`

func getCreateDataQuery(tableName string) string {
	numFields := dataFieldsCount[tableName]
	s := strings.Replace(queryCreateData, "{table}", tableName, -1)
	cont := ""
	for i := 0; i < numFields; i++ {
		cont += fmt.Sprintf("content_%d bytea DEFAULT NULL,\n", i+1)
	}
	s = strings.Replace(s, "{content}", cont, -1)
	return s
}

var queryCreateMetaData string = `CREATE TABLE IF NOT EXISTS public.metadata_{table}
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    data_id uuid NOT NULL REFERENCES public.data_{table} (id) ON DELETE CASCADE,
    name text NOT NULL,
    content bytea NOT NULL,
    md5 char(32),
    PRIMARY KEY (id),
    CONSTRAINT uk_metadata_{table}_name UNIQUE (data_id, name)
)
WITH (
    OIDS = FALSE
);`

func getCreateMetaDataQuery(tableName string) string {
	return strings.Replace(queryCreateMetaData, "{table}", tableName, -1)
}
