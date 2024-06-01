package main

import (
	"context"
	"go.uber.org/zap"
	"passwordvault/internal/config"
	"passwordvault/internal/grpc_server"
	"passwordvault/internal/storage"
	"time"
)

func main() {
	parentCtx, _ := context.WithCancel(context.Background())
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal(err.Error())
	}
	db, err := storage.New(config.GetServerConfig(), logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = db.Init(parentCtx)
	if err != nil {
		logger.Error(err.Error())
	}
	srv := grpc_server.NewGRPCServer(db, config.GetServerConfig(), logger)
	srv.ListenAndServeAsync()

	for {
		time.Sleep(5 * time.Second)
	}
}
