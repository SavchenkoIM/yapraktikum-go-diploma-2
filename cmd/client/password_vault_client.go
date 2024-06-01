package main

import (
	"context"
	"go.uber.org/zap"
	"passwordvault/internal/config"
	"passwordvault/internal/grpc_client"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	client := grpc_client.NewGRPCClient(config.GetClientConfig(), logger)

	client.Start(ctx)

	userData, err := client.UserLogin(ctx, "ivansav", "pass")
	if err != nil {
		panic(err)
	}
	token := userData.AccessToken

	err = client.PrintAllData(ctx, token)
	if err != nil {
		panic(err)
	}

	cancel()
}
