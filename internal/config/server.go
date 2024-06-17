package config

import (
	"errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var srvOnce sync.Once
var srvConfig *ServerConfig

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

func GetServerConfig() (*ServerConfig, error) {
	err := make([]error, 0)
	srvOnce.Do(func() {
		srvConfig = &ServerConfig{}

		flags := pflag.NewFlagSet("server", pflag.ExitOnError)
		flags.String("ag", ":8081", "gRPC Server endpoint address:port")
		flags.String("ah", ":8080", "HTTP Server endpoint address:port")
		flags.String("am", "localhost:9000", "MinIo Server endpoint address:port")
		flags.String("mcid", "minioadmin", "MinIo Administrator User ID")
		flags.String("mckey", "minioadmin", "MinIo Administrator Access Key")
		flags.String("d", "postgresql://localhost:5432/postgres?user=postgres&password=postgres", "Database connection string")
		flags.String("key", "", "Data encryption key")
		flags.String("cf", "", "Certificate file")
		flags.String("pkf", "", "Private key file")

		err = append(err, flags.Parse(os.Args))

		err = append(err, viper.BindPFlags(flags))

		err = append(err, viper.BindEnv("ag", "GRPC_ADDRESS"))
		err = append(err, viper.BindEnv("ah", "HTTP_ADDRESS"))
		err = append(err, viper.BindEnv("am", "MINIO_ADDRESS"))
		err = append(err, viper.BindEnv("mcid", "MINIO_ADMIN_ID"))
		err = append(err, viper.BindEnv("mckey", "MINIO_ADMIN_KEY"))
		err = append(err, viper.BindEnv("d", "DATABASE"))
		err = append(err, viper.BindEnv("key", "KEY"))
		err = append(err, viper.BindEnv("cf", "CERTIFICATE_FILENAME"))
		err = append(err, viper.BindEnv("pkf", "PK_FILENAME"))

		srvConfig.GrpcEndPoint = viper.GetString("ag")
		srvConfig.HttpEndPoint = viper.GetString("ah")
		srvConfig.MinioEndPoint = viper.GetString("mcid")
		srvConfig.MinioAdminId = viper.GetString("mckey")
		srvConfig.MinioAdminKey = viper.GetString("mckey")
		srvConfig.DBConnectionString = viper.GetString("d")
		srvConfig.Key = viper.GetString("key")
		srvConfig.UseKey = true
		srvConfig.CertFileName = viper.GetString("cf")
		srvConfig.PKFileName = viper.GetString("pkf")
	})
	return srvConfig, errors.Join(err...)
}
