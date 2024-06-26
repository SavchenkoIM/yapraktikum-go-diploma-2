package cli

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// User handler function delegate
type userHandleFunc func(context.Context, string, string) (string, error)

// Root "user" command
func (cli *CliManager) userCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	c.PersistentFlags().StringP("login", "l", "", "User Login")
	c.PersistentFlags().StringP("password", "p", "", "User Password")

	c.AddCommand(cli.userLogin())
	c.AddCommand(cli.userCreate())
	return c
}

// Unified action for user login and create
func (cli *CliManager) abstractUserAction(cmd *cobra.Command, act userHandleFunc) error {
	l, _ := cmd.Flags().GetString("login")
	p, _ := cmd.Flags().GetString("password")

	token, err := act(cmd.Context(), l, p)
	if err != nil {
		return err
	}
	viper.Set("token", token)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

// User login command
func (cli *CliManager) userLogin() *cobra.Command {
	c := &cobra.Command{
		Use:   "login",
		Short: "User login",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.abstractUserAction(cmd, cli.client.UserLogin)
		},
	}
	return c
}

// User create command
func (cli *CliManager) userCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "User registration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.abstractUserAction(cmd, cli.client.UserCreate)
		},
	}
	return c
}
