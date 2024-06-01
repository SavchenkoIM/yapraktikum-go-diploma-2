package config

import (
	"flag"
	"github.com/ianschenck/envflag"
	"slices"
	"sync"
)

var srvOnce sync.Once
var srvConfig *ServerConfig

type serverConfigNull struct {
	EndPoint           *string
	DBConnectionString *string
	Key                *string
	CertFileName       *string
	PKFileName         *string
}

type ServerConfig struct {
	EndPoint           string
	DBConnectionString string
	Key                string
	UseKey             bool
	CertFileName       string
	PKFileName         string
}

func GetServerConfig() *ServerConfig {
	srvOnce.Do(func() {
		srvConfig = CombineServerConfigs(getServerConfigFromCLArgs(), getServerConfigFromEnvVar())
	})
	return srvConfig
}

// Parses Server configuration from Command Line args
func getServerConfigFromCLArgs() serverConfigNull {
	serverConfig := serverConfigNull{}
	endPoint := flag.String("a", "", "Server endpoint address:port")
	dbConnString := flag.String("d", "", "Database connection string")
	encKey := flag.String("key", "", "Data encryption key")
	certFile := flag.String("cf", "", "Certificate file")
	pkFile := flag.String("pkf", "", "Private key file")
	flag.Parse()

	usedFlags := getProvidedFlags(flag.Visit)

	serverConfig.EndPoint = getParWithSetCheck(*endPoint, slices.Contains(usedFlags, "a"))
	serverConfig.DBConnectionString = getParWithSetCheck(*dbConnString, slices.Contains(usedFlags, "d"))
	serverConfig.Key = getParWithSetCheck(*encKey, slices.Contains(usedFlags, "key"))
	serverConfig.CertFileName = getParWithSetCheck(*certFile, slices.Contains(usedFlags, "cf"))
	serverConfig.PKFileName = getParWithSetCheck(*pkFile, slices.Contains(usedFlags, "pkf"))

	return serverConfig
}

// Parses Server configuration from Enviroment Vars
func getServerConfigFromEnvVar() serverConfigNull {
	serverConfig := serverConfigNull{}
	endPoint := envflag.String("ADDRESS", "", "Server endpoint address:port")
	dbConnString := envflag.String("CONN_STRING", "", "Database connection string")
	encKey := envflag.String("KEY", "", "Data encryption key")
	certFile := flag.String("CERT_FILE", "", "Certificate file")
	pkFile := flag.String("PK_FILE", "", "Private key file")
	envflag.Parse()

	usedFlags := getProvidedFlags(envflag.Visit)

	serverConfig.EndPoint = getParWithSetCheck(*endPoint, slices.Contains(usedFlags, "ADDRESS"))
	serverConfig.DBConnectionString = getParWithSetCheck(*dbConnString, slices.Contains(usedFlags, "CONN_STRING"))
	serverConfig.Key = getParWithSetCheck(*encKey, slices.Contains(usedFlags, "KEY"))
	serverConfig.CertFileName = getParWithSetCheck(*certFile, slices.Contains(usedFlags, "CERT_FILE"))
	serverConfig.PKFileName = getParWithSetCheck(*pkFile, slices.Contains(usedFlags, "PK_FILE"))

	return serverConfig
}

func CombineServerConfigs(configs ...serverConfigNull) *ServerConfig {
	serverConfig := ServerConfig{
		EndPoint:           ":8080",
		DBConnectionString: "postgresql://localhost:5432/postgres?user=postgres&password=postgres",
		Key:                "",
		UseKey:             false,
		CertFileName:       "",
		PKFileName:         "",
	}

	slices.Reverse(configs)
	for _, cfg := range configs {
		combineParameter(&serverConfig.EndPoint, cfg.EndPoint)
		combineParameter(&serverConfig.DBConnectionString, cfg.DBConnectionString)
		combineParameter(&serverConfig.Key, cfg.Key)
		combineParameter(&serverConfig.CertFileName, cfg.CertFileName)
		combineParameter(&serverConfig.PKFileName, cfg.PKFileName)
		if cfg.Key != nil && *cfg.Key != "" {
			serverConfig.UseKey = true
		}
	}

	return &serverConfig
}
