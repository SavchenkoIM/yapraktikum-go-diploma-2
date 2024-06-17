package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"passwordvault/internal/config"
	"passwordvault/internal/globals"
	"passwordvault/internal/grpc_server"
	"passwordvault/internal/http_server"
	"passwordvault/internal/storage/server_store"
	"syscall"
)

func main() {
	parentCtx, cancel := context.WithCancel(context.Background())
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Sugar().Infof("Password Vault Server. Version: v.%s (%s)", globals.ServerVer, globals.ServerDate)

	cfg, err := config.GetServerConfig()
	if err != nil {
		logger.Warn(err.Error())
	}

	db, err := server_store.New(cfg, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = db.Init(parentCtx)
	// CR#1: Not fatal. Storage will reconnect automatically after Postgres became reachable
	if err != nil {
		logger.Error(err.Error())
	}

	gSrv := grpc_server.NewGRPCServer(db, cfg, logger)
	gSrv.ListenAndServeAsync()

	hSrv := http_server.NewHttpServer(parentCtx, db, cfg, logger)
	hSrv.ListenAndServeAsync()

	gracefulShutdown(parentCtx, cancel, logger, gSrv, db)
}

func gracefulShutdown(
	ctx context.Context,
	cancel context.CancelFunc,
	logger *zap.Logger,
	gSrv *grpc_server.GRPCServer,
	db *server_store.Storage) {
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGTERM, syscall.SIGINT)
	s := <-terminateSignals
	logger.Info("Got one of stop signals, shutting down server gracefully, SIGNAL NAME :" + s.String())
	cancel()
	err := gSrv.Shutdown(ctx)
	if err != nil {
		logger.Error(err.Error())
	}
	db.Close(ctx)
}
