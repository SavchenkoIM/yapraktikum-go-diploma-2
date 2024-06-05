package grpc_client

import (
	"context"
	"crypto/tls"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
)

type GRPCClient struct {
	conn      *grpc.ClientConn
	cfg       *config.ClientConfig
	logger    *zap.Logger
	client    proto.PasswordVaultServiceClient
	token     string
	sendError chan error
}

func NewGRPCClient(cfg *config.ClientConfig, logger *zap.Logger) *GRPCClient {
	return &GRPCClient{cfg: cfg, logger: logger, sendError: make(chan error, 1)}
}

func (c *GRPCClient) Start(ctx context.Context) {
	// Connection
	var err error
	c.conn, err = grpc.NewClient(c.cfg.AddressGRPC, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		ClientAuth:         tls.NoClientCert, // Server provides cert
		InsecureSkipVerify: true,             // Any server_store cert is accepted
	})), grpc.WithUnaryInterceptor(c.WithUserCredentials))
	if err != nil {
		c.logger.Error("Failed to connect to server_store", zap.Error(err))
	}
	c.client = proto.NewPasswordVaultServiceClient(c.conn)
}

func (c *GRPCClient) Stop(ctx context.Context) error {
	return c.conn.Close()
}

func (c *GRPCClient) UserLogin(ctx context.Context, username string, password string) (*proto.UserResponse, error) {
	userRes, err := c.client.UserLogin(ctx, &proto.UserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	c.token = userRes.AccessToken
	return userRes, nil
}

func (c *GRPCClient) UserCreate(ctx context.Context, username string, password string) (*proto.UserResponse, error) {
	userRes, err := c.client.UserCreate(ctx, &proto.UserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	c.token = userRes.AccessToken
	return userRes, nil
}

func (c *GRPCClient) SetToken(token string) {
	c.token = token
}

func (c *GRPCClient) DataWrite(ctx context.Context, request *proto.DataWriteRequest) (*proto.EmptyResponse, error) {
	return c.client.DataWrite(ctx, request)
}

func (c *GRPCClient) DataRead(ctx context.Context, filter *proto.DataReadRequest) (*proto.DataReadResponse, error) {
	return c.client.DataRead(ctx, filter)
}
