package main

import (
	"context"
	"fmt"
	"passwordvault/internal/uni_client/cli"
)

/*
Will print:

Credentials:
	Name: my_main_creds
	Login: victoria
	Password: victoria's secret
	Metadata:
		site: google.com
==================
Files:
	Name: my_first_script
	FileName: test.py
	Metadata:
		code_quality: bad
==================

And will download 'test.py' to default folder on client PC
*/

func main() {

	err := (&cli.CliManager{}).ExecuteContext(context.Background())
	if err != nil {
		fmt.Println(err)
	}

}
