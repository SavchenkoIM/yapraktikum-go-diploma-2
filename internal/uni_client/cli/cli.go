// Clommand line interface for client

package cli

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"passwordvault/internal/config"
	"passwordvault/internal/uni_client"
	"path/filepath"
)

// Object, containing global data for cli logic
type CliManager struct {
	client *uni_client.UniClient
	Logger *zap.Logger
}

// Execution ot root command
func (cli *CliManager) ExecuteContext(ctx context.Context) error {
	return cli.rootCmd().ExecuteContext(ctx)
}

// Root command
func (cli *CliManager) rootCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "password_vault_client",
	}

	c.PersistentFlags().StringP("grpc-address", "g", "localhost:8081", "gRPC Server address:port")
	c.PersistentFlags().StringP("http-address", "a", "localhost:8080", "HTTP Server address:port")
	c.PersistentFlags().StringP("files-dir", "f", "", "Files default directory")
	c.PersistentFlags().StringP("config", "c", "", "Config file path")

	viper.BindPFlag("grpc-address", c.PersistentFlags().Lookup("grpc-address"))
	viper.BindPFlag("http-address", c.PersistentFlags().Lookup("http-address"))
	viper.BindPFlag("files-dir", c.PersistentFlags().Lookup("files-dir"))
	viper.BindPFlag("config", c.PersistentFlags().Lookup("config"))

	viper.BindEnv("grpc-address", "GRPC_ADDRESS",
		"http-address", "HTTP_ADDRESS",
		"files-dir", "FILES_DIR",
		"config", "CONFIG_FILE")

	cli.initConfig(viper.GetString("config"))

	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		cli.client = uni_client.NewUniClient(cli.Logger, config.ClientConfig{
			AddressGRPC:     viper.GetString("grpc-address"),
			AddressHTTP:     viper.GetString("http-address"),
			FilesDefaultDir: viper.GetString("files-dir"),
		})
		cli.client.Start(cmd.Context())

		token := viper.GetString("token")
		cli.client.SetToken(token)

		return nil
	}

	c.AddCommand(cli.dataCommand())
	c.AddCommand(cli.fileCommand())
	c.AddCommand(cli.userCommand())

	return c
}

// Reads client configuration from different sources: cl args, envs and config file
func (cli *CliManager) initConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			cli.Logger.Fatal(err.Error())
		}
		viper.SetConfigFile(filepath.Join(home, "pass_vault.yaml"))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if cfgFile != "" {
			cli.Logger.Info("config specified but unable to read it, using defaults")
		}
	}
}
