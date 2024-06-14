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
	GrpcEndPoint       *string
	HttpEndPoint       *string
	MinioEndPoint      *string
	MinioAdminId       *string
	MinioAdminKey      *string
	DBConnectionString *string
	Key                *string
	CertFileName       *string
	PKFileName         *string
}

type ServerConfig struct {
	GrpcEndPoint       string
	HttpEndPoint       string
	MinioEndPoint      string
	MinioAdminId       string
	MinioAdminKey      string
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
	endPoint := flag.String("ag", "", "gRPC Server endpoint address:port")
	httpEndPoint := flag.String("ah", "", "HTTP Server endpoint address:port")
	minioEndPoint := flag.String("am", "", "MinIo Server endpoint address:port")
	minioAdminId := flag.String("mcid", "", "MinIo Administrator User ID")
	minioAdminKey := flag.String("mckey", "", "MinIo Administrator Access Key")
	dbConnString := flag.String("d", "", "Database connection string")
	encKey := flag.String("key", "", "Data encryption key")
	certFile := flag.String("cf", "", "Certificate file")
	pkFile := flag.String("pkf", "", "Private key file")
	flag.Parse()

	usedFlags := getProvidedFlags(flag.Visit)

	serverConfig.GrpcEndPoint = getParWithSetCheck(*endPoint, slices.Contains(usedFlags, "ag"))
	serverConfig.HttpEndPoint = getParWithSetCheck(*httpEndPoint, slices.Contains(usedFlags, "ah"))
	serverConfig.MinioEndPoint = getParWithSetCheck(*minioEndPoint, slices.Contains(usedFlags, "am"))
	serverConfig.MinioAdminId = getParWithSetCheck(*minioAdminId, slices.Contains(usedFlags, "mcid"))
	serverConfig.MinioAdminKey = getParWithSetCheck(*minioAdminKey, slices.Contains(usedFlags, "mckey"))
	serverConfig.DBConnectionString = getParWithSetCheck(*dbConnString, slices.Contains(usedFlags, "d"))
	serverConfig.Key = getParWithSetCheck(*encKey, slices.Contains(usedFlags, "key"))
	serverConfig.CertFileName = getParWithSetCheck(*certFile, slices.Contains(usedFlags, "cf"))
	serverConfig.PKFileName = getParWithSetCheck(*pkFile, slices.Contains(usedFlags, "pkf"))

	return serverConfig
}

// Parses Server configuration from Enviroment Vars
func getServerConfigFromEnvVar() serverConfigNull {
	serverConfig := serverConfigNull{}
	endPoint := envflag.String("GRPC_ADDRESS", "", "gRPC Server endpoint address:port")
	httpEndPoint := envflag.String("HTTP_ADDRESS", "", "HTTP Server endpoint address:port")
	minioEndPoint := envflag.String("MINIO_ADDRESS", "", "MinIo Server endpoint address:port")
	minioAdminId := envflag.String("MINIO_ADMIN_ID", "", "MinIo Administrator User ID")
	minioAdminKey := envflag.String("MINIO_ADMIN_KEY", "", "MinIo Administrator Access Key")
	dbConnString := envflag.String("CONN_STRING", "", "Database connection string")
	encKey := envflag.String("KEY", "", "Data encryption key")
	certFile := envflag.String("CERT_FILE", "", "Certificate file")
	pkFile := envflag.String("PK_FILE", "", "Private key file")
	envflag.Parse()

	usedFlags := getProvidedFlags(envflag.Visit)

	serverConfig.GrpcEndPoint = getParWithSetCheck(*endPoint, slices.Contains(usedFlags, "GRPC_ADDRESS"))
	serverConfig.HttpEndPoint = getParWithSetCheck(*httpEndPoint, slices.Contains(usedFlags, "HTTP_ADDRESS"))
	serverConfig.MinioEndPoint = getParWithSetCheck(*minioEndPoint, slices.Contains(usedFlags, "MINIO_ADDRESS"))
	serverConfig.MinioAdminId = getParWithSetCheck(*minioAdminId, slices.Contains(usedFlags, "MINIO_ADMIN_ID"))
	serverConfig.MinioAdminKey = getParWithSetCheck(*minioAdminKey, slices.Contains(usedFlags, "MINIO_ADMIN_KEY"))
	serverConfig.DBConnectionString = getParWithSetCheck(*dbConnString, slices.Contains(usedFlags, "CONN_STRING"))
	serverConfig.Key = getParWithSetCheck(*encKey, slices.Contains(usedFlags, "KEY"))
	serverConfig.CertFileName = getParWithSetCheck(*certFile, slices.Contains(usedFlags, "CERT_FILE"))
	serverConfig.PKFileName = getParWithSetCheck(*pkFile, slices.Contains(usedFlags, "PK_FILE"))

	return serverConfig
}

// Concatenates server configs from different sources
func CombineServerConfigs(configs ...serverConfigNull) *ServerConfig {
	serverConfig := ServerConfig{
		GrpcEndPoint:       ":8081",
		HttpEndPoint:       ":8080",
		MinioEndPoint:      "localhost:9000",
		MinioAdminId:       "minioadmin",
		MinioAdminKey:      "minioadmin",
		DBConnectionString: "postgresql://localhost:5432/postgres?user=postgres&password=postgres",
		Key:                "",
		UseKey:             false,
		CertFileName:       "",
		PKFileName:         "",
	}

	slices.Reverse(configs)
	for _, cfg := range configs {
		combineParameter(&serverConfig.GrpcEndPoint, cfg.GrpcEndPoint)
		combineParameter(&serverConfig.HttpEndPoint, cfg.HttpEndPoint)
		combineParameter(&serverConfig.MinioEndPoint, cfg.MinioEndPoint)
		combineParameter(&serverConfig.MinioAdminId, cfg.MinioAdminId)
		combineParameter(&serverConfig.MinioAdminKey, cfg.MinioAdminKey)
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
