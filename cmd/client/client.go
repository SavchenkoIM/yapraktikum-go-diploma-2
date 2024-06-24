package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"passwordvault/internal/globals"
	"passwordvault/internal/uni_client/cli"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger.Sugar().Infof("Password Vault Client. Version: v.%s (%s)", globals.ClientVer, globals.ClientDate)
	if err := (&cli.CliManager{Logger: logger}).ExecuteContext(context.Background()); err != nil {
		fmt.Println(err)
	}
}
