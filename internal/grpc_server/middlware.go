package grpc_server

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"passwordvault/internal/config"
	"strings"
	"time"
)

// Middlware for gRPC server
type GRPCServerMiddleware struct {
	Cfg    config.ServerConfig
	Logger *zap.Logger
}

// Logging middlware for gRPC server
func (s *GRPCServer) WithLogging(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Calls the handler
	h, err := handler(ctx, req)

	// Logging with grpclog (grpclog.LoggerV2)
	s.logger.Sugar().Infof("gRPC Method: %+v, Runtime: %d msec, Error:%v",
		info.FullMethod, time.Since(start).Milliseconds(), err)

	return h, err
}

// Middlware handler for authentication data for gRPC server
func (s *GRPCServer) authFunc(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	if strings.HasPrefix(method, "/grpc.PasswordVaultService/User") {
		return ctx, nil
	}

	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	login, err := s.dataStorage.UserCheckLoggedIn(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user auth error: %v", err)
	}

	return context.WithValue(ctx, "LoggedUserId", login), nil
}
