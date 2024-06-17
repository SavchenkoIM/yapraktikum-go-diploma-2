// gRPC server implementation

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
	"passwordvault/internal/storage/server_store"
	"passwordvault/internal/utils"
)

// gRPC server object data
type GRPCServer struct {
	proto.UnsafePasswordVaultServiceServer
	dataStorage *server_store.Storage
	gsrv        *grpc.Server
	srv         net.Listener
	cfg         *config.ServerConfig
	logger      *zap.Logger
}

// Constructor of gRPC server object
func NewGRPCServer(dataStorage *server_store.Storage, cfg *config.ServerConfig, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{dataStorage: dataStorage, cfg: cfg, logger: logger}
}

// Starts gRPC server listener asyncronously
func (s *GRPCServer) ListenAndServeAsync() error {
	var err error
	s.srv, err = net.Listen("tcp", s.cfg.GrpcEndPoint)
	if err != nil {
		return err
	}

	creds, err := utils.LoadTLSCredentials(s.cfg.CertFileName, s.cfg.PKFileName)
	if err != nil {
		return err
	}
	s.gsrv = grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			s.WithLogging,
			auth.UnaryServerInterceptor(s.authFunc),
			recovery.UnaryServerInterceptor()))
	proto.RegisterPasswordVaultServiceServer(s.gsrv, s)

	go func() {
		s.logger.Info("gRPC server running at " + s.cfg.GrpcEndPoint)

		if err := s.gsrv.Serve(s.srv); err != nil {
			s.logger.Error(err.Error())
			return
		}
	}()

	return nil
}

// Stops server
func (s *GRPCServer) Shutdown(context.Context) error {
	s.gsrv.Stop()
	return nil
}
