package http_server

import (
	"context"
	"crypto/tls"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/storage/server_store"
)

type HttpServer struct {
	server *http.Server
	logger *zap.Logger
	config *config.ServerConfig
	db     *server_store.Storage
}

func NewHttpServer(ctx context.Context, db *server_store.Storage, config *config.ServerConfig, logger *zap.Logger) *HttpServer {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		ClientAuth:         tls.NoClientCert, // Server provides cert
		InsecureSkipVerify: true,             // Any server_store cert is accepted
	}))}

	err := proto.RegisterPasswordVaultServiceHandlerFromEndpoint(ctx, mux, config.GrpcEndPoint, opts)
	if err != nil {
		logger.Fatal(err.Error())
	}

	srv := &HttpServer{
		server: &http.Server{Addr: config.HttpEndPoint, Handler: mux},
		logger: logger,
		config: config,
		db:     db,
	}

	err = mux.HandlePath("POST", `/download`, srv.WithLoggingHTTP(srv.WithCheckCredentials(srv.DownloadFile)))
	if err != nil {
		srv.logger.Fatal(err.Error())
	}
	err = mux.HandlePath("POST", `/upload`, srv.WithLoggingHTTP(srv.WithCheckCredentials(srv.UploadFile)))
	if err != nil {
		srv.logger.Fatal(err.Error())
	}

	return srv
}

func (s *HttpServer) ListenAndServeAsync() {
	go func() {
		s.logger.Info("HTTP server running at " + s.config.HttpEndPoint)

		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Error(err.Error())
			return
		}
	}()
}
