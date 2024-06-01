package grpc_server

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/storage"
	"passwordvault/internal/utils"
)

type GRPCServer struct {
	proto.UnsafePasswordVaultServiceServer
	dataStorage *storage.Storage
	gsrv        *grpc.Server
	srv         net.Listener
	cfg         *config.ServerConfig
	logger      *zap.Logger
}

func NewGRPCServer(dataStorage *storage.Storage, cfg *config.ServerConfig, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{dataStorage: dataStorage, cfg: cfg, logger: logger}
}

func (s *GRPCServer) ListenAndServeAsync() {
	var err error
	s.srv, err = net.Listen("tcp", s.cfg.EndPoint)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	creds, err := utils.LoadTLSCredentials(s.cfg.CertFileName, s.cfg.PKFileName)
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	s.gsrv = grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			s.WithLogging,
			auth.UnaryServerInterceptor(s.authFunc),
			recovery.UnaryServerInterceptor()))
	proto.RegisterPasswordVaultServiceServer(s.gsrv, s)

	go func() {
		s.logger.Info("gRPC server running at " + s.cfg.EndPoint)

		if err := s.gsrv.Serve(s.srv); err != nil {
			s.logger.Error(err.Error())
			return
		}
	}()
}

func (s *GRPCServer) Shutdown(context.Context) error {
	s.gsrv.Stop()
	return nil
}
