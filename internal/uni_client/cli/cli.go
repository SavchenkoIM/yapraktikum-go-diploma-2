package cli

import (
	"context"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"passwordvault/internal/config"
	"passwordvault/internal/uni_client"
)

type CliManager struct {
	client *uni_client.UniClient
	logger *zap.Logger
}

func (cli *CliManager) ExecuteContext(ctx context.Context) error {
	var err error
	cli.logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return cli.rootCmd().ExecuteContext(ctx)
}

func (cli *CliManager) rootCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "password_vault_client",
	}

	c.PersistentFlags().StringP("grpc-address", "g", "localhost:8081", "gRPC Server address:port")
	c.PersistentFlags().StringP("http-address", "a", "localhost:8080", "HTTP Server address:port")
	c.PersistentFlags().StringP("files-dir", "f", "", "Files default directory")
	c.PersistentFlags().StringP("login", "l", "", "User Login")
	c.PersistentFlags().StringP("password", "p", "", "User Password")

	c.AddCommand(cli.dataCommand())
	c.AddCommand(cli.fileCommand())

	return c
}

func (cli *CliManager) initClient() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		l, _ := cmd.Flags().GetString("login")
		p, _ := cmd.Flags().GetString("password")
		ag, _ := cmd.Flags().GetString("grpc-address")
		ah, _ := cmd.Flags().GetString("http-address")
		f, _ := cmd.Flags().GetString("files-dir")

		cli.client = uni_client.NewUniClient(cli.logger, config.ClientConfig{
			AddressGRPC:     ag,
			AddressHTTP:     ah,
			FilesDefaultDir: f,
		})
		cli.client.Start(cmd.Context())

		_, err := cli.client.UserLogin(cmd.Context(), l, p)
		return err
	}
}
