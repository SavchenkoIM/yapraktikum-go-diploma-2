package config

import (
	"flag"
	"github.com/ianschenck/envflag"
	"slices"
	"sync"
)

var cliOnce sync.Once
var cliConfig *ClientConfig

type clientConfigNull struct {
	AddressGRPC     *string
	AddressHTTP     *string
	FilesDefaultDir *string
}

type ClientConfig struct {
	AddressGRPC     string
	AddressHTTP     string
	FilesDefaultDir string
}

/*func GetClientConfig() *ClientConfig {
	cliOnce.Do(func() {
		cliConfig = CombineClientConfigs(getClientConfigFromCLArgs(), getClientConfigFromEnvVar())
	})
	return cliConfig
}*/

// Parses Client configuration from Command Line args
func getClientConfigFromCLArgs() clientConfigNull {
	clientConfig := clientConfigNull{}
	addressGrpc := flag.String("ag", "", "gRPC Server address:port")
	addressHttp := flag.String("ah", "", "HTTP Server address:port")
	filesDefaultDir := flag.String("f", "", "Files default directory")
	flag.Parse()

	usedFlags := getProvidedFlags(flag.Visit)

	clientConfig.AddressGRPC = getParWithSetCheck(*addressGrpc, slices.Contains(usedFlags, "ag"))
	clientConfig.AddressHTTP = getParWithSetCheck(*addressHttp, slices.Contains(usedFlags, "ah"))
	clientConfig.FilesDefaultDir = getParWithSetCheck(*filesDefaultDir, slices.Contains(usedFlags, "f"))

	return clientConfig
}

// Parses Client configuration from Enviroment Vars
func getClientConfigFromEnvVar() clientConfigNull {
	clientConfig := clientConfigNull{}
	addressGrpc := envflag.String("ADDRESS_GRPC", "", "gRPC Server address:port")
	addressHttp := envflag.String("ADDRESS_HTTP", "", "HTTP Server address:port")
	filesDefaultDir := envflag.String("FILES_DEFAULT_DIR", "", "Files default directory")
	envflag.Parse()

	usedFlags := getProvidedFlags(envflag.Visit)

	clientConfig.AddressGRPC = getParWithSetCheck(*addressGrpc, slices.Contains(usedFlags, "ADDRESS_GRPC"))
	clientConfig.AddressHTTP = getParWithSetCheck(*addressHttp, slices.Contains(usedFlags, "ADDRESS_HTTP"))
	clientConfig.FilesDefaultDir = getParWithSetCheck(*filesDefaultDir, slices.Contains(usedFlags, "FILES_DEFAULT_DIR"))

	return clientConfig
}

func CombineClientConfigs(configs ...clientConfigNull) *ClientConfig {
	clientConfig := ClientConfig{
		AddressGRPC:     ":8081",
		AddressHTTP:     ":8080",
		FilesDefaultDir: ".",
	}

	slices.Reverse(configs)
	for _, cfg := range configs {
		combineParameter(&clientConfig.AddressGRPC, cfg.AddressGRPC)
		combineParameter(&clientConfig.AddressHTTP, cfg.AddressHTTP)
		combineParameter(&clientConfig.FilesDefaultDir, cfg.FilesDefaultDir)
	}

	return &clientConfig
}
