package main

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"passwordvault/internal/config"
	"passwordvault/internal/grpc_server"
	"passwordvault/internal/http_server"
	"passwordvault/internal/storage/server_store"
	"time"
)

var (
	BuildDate    string
	BuildVersion string
)

func WithUserCredentialsFrom(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return handler(ctx, req)
}

func main() {
	parentCtx, _ := context.WithCancel(context.Background())
	logger, err := zap.NewProduction()

	logger.Sugar().Infof("Password Vault Server. Version: v.%s (%s)", BuildVersion, BuildDate)

	if err != nil {
		logger.Fatal(err.Error())
	}
	db, err := server_store.New(config.GetServerConfig(), logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = db.Init(parentCtx)
	if err != nil {
		logger.Error(err.Error())
	}

	gSrv := grpc_server.NewGRPCServer(db, config.GetServerConfig(), logger)
	gSrv.ListenAndServeAsync()

	hSrv := http_server.NewHttpServer(parentCtx, db, config.GetServerConfig(), logger)
	hSrv.ListenAndServeAsync()

	for {
		time.Sleep(5 * time.Second)
	}
}
