package cli

import (
	"github.com/spf13/cobra"
)

// Root "file" command
func (cli *CliManager) fileCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "file",
		Short: "Upload and download files",
	}

	c.AddCommand(cli.fileUpload())
	c.AddCommand(cli.fileDownload())
	return c
}

// File upload command
func (cli *CliManager) fileUpload() *cobra.Command {
	c := &cobra.Command{
		Use:   "upload",
		Short: "Upload file",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			f, _ := cmd.Flags().GetString("fname")
			return cli.client.UploadFile(cmd.Context(), n, f)
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")
	c.Flags().String("fname", "", "File name on disc")
	c.MarkFlagRequired("fname")

	return c
}

// File download command
func (cli *CliManager) fileDownload() *cobra.Command {
	c := &cobra.Command{
		Use:   "download",
		Short: "Download file",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := cmd.Flags().GetString("name")
			return cli.client.DownloadFile(cmd.Context(), n)
		},
	}

	c.Flags().StringP("name", "n", "", "Name of data record")
	c.MarkFlagRequired("name")

	return c
}
