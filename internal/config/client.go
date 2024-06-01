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
	Address *string
}

type ClientConfig struct {
	Address string
}

func GetClientConfig() *ClientConfig {
	cliOnce.Do(func() {
		cliConfig = CombineClientConfigs(getClientConfigFromCLArgs(), getClientConfigFromEnvVar())
	})
	return cliConfig
}

// Parses Client configuration from Command Line args
func getClientConfigFromCLArgs() clientConfigNull {
	clientConfig := clientConfigNull{}
	address := flag.String("a", "", "Server address:port")
	flag.Parse()

	usedFlags := getProvidedFlags(flag.Visit)

	clientConfig.Address = getParWithSetCheck(*address, slices.Contains(usedFlags, "a"))

	return clientConfig
}

// Parses Client configuration from Enviroment Vars
func getClientConfigFromEnvVar() clientConfigNull {
	clientConfig := clientConfigNull{}
	address := envflag.String("ADDRESS", "", "Server address:port")
	envflag.Parse()

	usedFlags := getProvidedFlags(envflag.Visit)

	clientConfig.Address = getParWithSetCheck(*address, slices.Contains(usedFlags, "ADDRESS"))

	return clientConfig
}

func CombineClientConfigs(configs ...clientConfigNull) *ClientConfig {
	clientConfig := ClientConfig{
		Address: ":8080",
	}

	slices.Reverse(configs)
	for _, cfg := range configs {
		combineParameter(&clientConfig.Address, cfg.Address)
	}

	return &clientConfig
}
