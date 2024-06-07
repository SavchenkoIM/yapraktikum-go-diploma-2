package cli

import (
	"github.com/spf13/cobra"
	proto "passwordvault/internal/proto/gen"
)

func (cli *CliManager) writeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "write",
		Short: "Write data or metadata record",
	}

	c.AddCommand(cli.writeCredentialsCommand())
	c.AddCommand(cli.writeCreditCardCommand())
	c.AddCommand(cli.writeTextNoteCommand())
	c.AddCommand(cli.writeMetadataCommand())

	return c
}

func (cli *CliManager) writeCredentialsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "cred",
		Short:   "Write credentials data record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			l, _ := cmd.Flags().GetString("dlogin")
			p, _ := cmd.Flags().GetString("dpassword")
			_, err := cli.client.DataWrite(cmd.Context(), &proto.DataWriteRequest{
				Action: proto.OperationType_UPSERT,
				Data: &proto.DataWriteRequest_Credentials{Credentials: &proto.DataCredentials{
					Name:     n,
					Login:    l,
					Password: p,
				}},
			})
			return err
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")
	c.Flags().String("dlogin", "", "Login")
	c.Flags().String("dpassword", "", "Password")
	c.MarkFlagRequired("dlogin")
	c.MarkFlagRequired("dpassword")

	return c
}

func (cli *CliManager) writeCreditCardCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "card",
		Short:   "Write credit card data record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			nu, _ := cmd.Flags().GetString("dnumber")
			u, _ := cmd.Flags().GetString("duntil")
			h, _ := cmd.Flags().GetString("dholder")
			_, err := cli.client.DataWrite(cmd.Context(), &proto.DataWriteRequest{
				Action: proto.OperationType_UPSERT,
				Data: &proto.DataWriteRequest_CreditCard{CreditCard: &proto.DataCreditCard{
					Name:   n,
					Number: nu,
					Until:  u,
					Holder: h,
				}},
			})
			return err
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")
	c.Flags().String("dnumber", "", "Card number")
	c.Flags().String("duntil", "", "Expiration date")
	c.Flags().String("dholder", "", "Holder name")
	c.MarkFlagRequired("dnumber")
	c.MarkFlagRequired("duntil")
	c.MarkFlagRequired("dholder")

	return c
}

func (cli *CliManager) writeTextNoteCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "note",
		Short:   "Write text note data record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			t, _ := cmd.Flags().GetString("dtext")
			_, err := cli.client.DataWrite(cmd.Context(), &proto.DataWriteRequest{
				Action: proto.OperationType_UPSERT,
				Data: &proto.DataWriteRequest_TextNote{TextNote: &proto.DataTextNote{
					Name: n,
					Text: t,
				}},
			})
			return err
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")
	c.Flags().String("dtext", "", "Text of note")
	c.MarkFlagRequired("dtext")

	return c
}

func (cli *CliManager) writeMetadataCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "metadata",
		Short:   "Write metadata for record",
		PreRunE: cli.initClient(),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			t, _ := cmd.Flags().GetString("type")
			nm, _ := cmd.Flags().GetString("dname")
			v, _ := cmd.Flags().GetString("dvalue")
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
					Value:      v,
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
	c.Flags().String("dvalue", "", "Metadata value")
	c.MarkFlagRequired("dname")
	c.MarkFlagRequired("dvalue")

	return c
}
