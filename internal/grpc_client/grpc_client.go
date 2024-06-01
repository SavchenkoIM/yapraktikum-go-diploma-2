package grpc_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
)

type GRPCClient struct {
	conn      *grpc.ClientConn
	cfg       *config.ClientConfig
	logger    *zap.Logger
	client    proto.PasswordVaultServiceClient
	sendError chan error
}

func NewGRPCClient(cfg *config.ClientConfig, logger *zap.Logger) *GRPCClient {
	return &GRPCClient{cfg: cfg, logger: logger, sendError: make(chan error, 1)}
}

func (c *GRPCClient) Start(ctx context.Context) {
	// Connection
	var err error
	c.conn, err = grpc.NewClient(c.cfg.Address, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: true, // Any server cert is accepted
	})))
	if err != nil {
		c.logger.Error("Failed to connect to server", zap.Error(err))
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
	return userRes, nil
}

func (c *GRPCClient) PrintAllData(ctx context.Context, token string) error {

	ctxf := metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)

	dataRes, err := c.client.DataRead(ctxf, &proto.DataReadRequest{
		Type:     proto.DataType_UNSPECIFIED,
		NameMask: "%",
		Metadata: nil,
	})
	if err != nil {
		return err
	}

	for _, v := range dataRes.Data {
		switch vv := v.Data.(type) {
		case *proto.DataRecord_CreditCard:
			fmt.Printf("CreditCard:\n\tName: %s\n\tNumber: %d\n\tUntil: %d\n\tHolder: %s",
				vv.CreditCard.Name,
				vv.CreditCard.Number,
				vv.CreditCard.Until,
				vv.CreditCard.Holder)
		case *proto.DataRecord_TextNote:
			fmt.Printf("TextNote:\n\tName: %s\n\tText: %s",
				vv.TextNote.Name,
				vv.TextNote.Text)
		case *proto.DataRecord_Credentials:
			fmt.Printf("Credentials:\n\tName: %s\n\tLogin: %s\n\tPassword: %s",
				vv.Credentials.Name,
				vv.Credentials.Login,
				vv.Credentials.Password)
		}

		if len(v.Metadata) > 0 {
			fmt.Print("\n\tMetadata:\n")
		}
		for _, md := range v.Metadata {
			fmt.Printf("\t\t%s: %s", md.Name, md.Value)
		}

		fmt.Println("\n==================")
	}

	return nil
}
