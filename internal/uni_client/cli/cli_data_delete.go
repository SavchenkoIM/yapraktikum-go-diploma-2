package cli

import (
	"github.com/spf13/cobra"
	proto "passwordvault/internal/proto/gen"
)

func (cli *CliManager) deleteCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete",
		Short: "Remove data or metadata record",
	}

	c.AddCommand(cli.deleteCredentialsCommand())
	c.AddCommand(cli.deleteCreditCardCommand())
	c.AddCommand(cli.deleteTextNoteCommand())
	c.AddCommand(cli.deleteMetadataCommand())

	return c
}

func (cli *CliManager) deleteCredentialsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "cred",
		Short:   "Remove credentials data record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			_, err := cli.client.DataWrite(cmd.Context(), &proto.DataWriteRequest{
				Action: proto.OperationType_DELETE,
				Data: &proto.DataWriteRequest_Credentials{Credentials: &proto.DataCredentials{
					Name: n,
				}},
			})
			return err
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")

	return c
}

func (cli *CliManager) deleteCreditCardCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "card",
		Short:   "Remove credit card data record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			_, err := cli.client.DataWrite(cmd.Context(), &proto.DataWriteRequest{
				Action: proto.OperationType_DELETE,
				Data: &proto.DataWriteRequest_CreditCard{CreditCard: &proto.DataCreditCard{
					Name: n,
				}},
			})
			return err
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")
	return c
}

func (cli *CliManager) deleteTextNoteCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "note",
		Short:   "Remove text note data record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			_, err := cli.client.DataWrite(cmd.Context(), &proto.DataWriteRequest{
				Action: proto.OperationType_UPSERT,
				Data: &proto.DataWriteRequest_TextNote{TextNote: &proto.DataTextNote{
					Name: n,
				}},
			})
			return err
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")
	return c
}

func (cli *CliManager) deleteMetadataCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "metadata",
		Short:   "Remove metadata for record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			t, _ := cmd.Flags().GetString("type")
			nm, _ := cmd.Flags().GetString("dname")
			dt, err := RecordType(t).GetType()
			if err != nil {
				cli.logger.Error(err.Error())
				return err
			}
			_, err = cli.client.DataWrite(cmd.Context(), &proto.DataWriteRequest{
				Action: proto.OperationType_UPSERT,
				Data: &proto.DataWriteRequest_Metadata{Metadata: &proto.MetaDataKV{
					ParentType: dt,
					ParentName: n,
					Name:       nm,
				}},
			})
			return err
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.Flags().StringP("type", "t", "", "Type of data record")
	c.MarkFlagRequired("name")
	c.MarkFlagRequired("type")
	c.Flags().String("dname", "", "Metadata name")
	c.MarkFlagRequired("dname")

	return c
}
