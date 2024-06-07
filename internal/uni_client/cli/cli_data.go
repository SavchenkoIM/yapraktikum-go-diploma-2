package cli

import (
	"github.com/spf13/cobra"
	proto "passwordvault/internal/proto/gen"
)

func (cli *CliManager) dataCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "data",
		Short: "Tools for managing amd accessing data",
	}

	c.AddCommand(cli.printCommand())
	c.AddCommand(cli.writeCommand())
	c.AddCommand(cli.deleteCommand())

	return c
}

func (cli *CliManager) printCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "print",
		Short:   "Print required data objects",
		PreRunE: cli.initClient(),
		Run: func(cmd *cobra.Command, args []string) {
			n, _ := cmd.Flags().GetString("name")
			t, _ := cmd.Flags().GetString("type")
			m, _ := cmd.Flags().GetStringSlice("metadata")
			dt, err := RecordType(t).GetType()
			if err != nil {
				cli.logger.Error(err.Error())
				return
			}
			fl, err := MetadataFilters(m).GetFilters()
			if err != nil {
				cli.logger.Error(err.Error())
				return
			}
			cli.client.DataPrint(cmd.Context(), &proto.DataReadRequest{
				Type:     dt,
				NameMask: n,
				Metadata: fl,
			})
		},
	}

	c.Flags().StringP("name", "n", "%", "Name mask of data record (accept SQL LIKE syntax)")
	c.Flags().StringP("type", "t", "ANY", "Record type")
	c.Flags().StringSliceP("metadata", "m", []string{}, "Metadata based filters")

	return c
}
