// Implementation of unified gRPC/HTTP client. Uses HTTP for upload and download file, gRPC for other operations.

package uni_client

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/uni_client/grpc_client"
	"passwordvault/internal/uni_client/http_client"
	"path/filepath"
)

// UniClient data
type UniClient struct {
	gCli   *grpc_client.GRPCClient
	hCli   *http_client.HTTPClient
	logger *zap.Logger
}

// Constructs UniClient data object
func NewUniClient(logger *zap.Logger, clientConfig config.ClientConfig) *UniClient {
	return &UniClient{
		gCli:   grpc_client.NewGRPCClient(&clientConfig, logger),
		hCli:   http_client.NewHTTPClient(&clientConfig, logger),
		logger: logger,
	}
}

// Starts client
func (c *UniClient) Start(ctx context.Context) {
	c.gCli.Start(ctx)
}

// Shuts client down
func (c *UniClient) Stop(ctx context.Context) error {
	return c.gCli.Stop(ctx)
}

// User login
func (c *UniClient) UserLogin(ctx context.Context, username string, password string) (string, error) {
	userData, err := c.gCli.UserLogin(ctx, username, password)
	if err != nil {
		return "", err
	}
	c.hCli.SetToken(userData.AccessToken)
	return userData.AccessToken, nil
}

// User create
func (c *UniClient) UserCreate(ctx context.Context, username string, password string) (string, error) {
	userData, err := c.gCli.UserCreate(ctx, username, password)
	if err != nil {
		return "", err
	}
	c.hCli.SetToken(userData.AccessToken)
	return userData.AccessToken, nil
}

// Sets default user auth token
func (c *UniClient) SetToken(token string) {
	c.hCli.SetToken(token)
	c.gCli.SetToken(token)
}

// Download file
func (c *UniClient) DownloadFile(ctx context.Context, objectName string) error {
	data, err := c.gCli.DataRead(ctx, &proto.DataReadRequest{
		Type:     proto.DataType_BLOB,
		NameMask: objectName,
		Metadata: make([]*proto.MetaDataKV, 0),
	})
	if err != nil {
		return err
	}
	if len(data.Data) != 1 {
		return errors.New(fmt.Sprintf("file with objectName: %s not found", objectName))
	}
	fileName := data.Data[0].Data.(*proto.DataRecord_Blob).Blob.FileName

	err = c.hCli.DownloadFile(ctx, fileName)
	if err != nil {
		return err
	}
	return nil
}

// Delete file
func (c *UniClient) DeleteFile(ctx context.Context, objectName string) error {
	_, err := c.gCli.DataWrite(ctx, &proto.DataWriteRequest{
		Action: proto.OperationType_DELETE,
		Data: &proto.DataWriteRequest_Blob{
			Blob: &proto.DataBLOB{
				Name: objectName,
			}},
	})
	if err != nil {
		return err
	}
	return nil
}

// Upload file
func (c *UniClient) UploadFile(ctx context.Context, objectName string, fileName string) error {
	err := c.hCli.UploadFile(ctx, fileName)
	if err != nil {
		return err
	}
	_, err = c.gCli.DataWrite(ctx, &proto.DataWriteRequest{
		Action: proto.OperationType_UPSERT,
		Data: &proto.DataWriteRequest_Blob{
			Blob: &proto.DataBLOB{
				Name:     objectName,
				FileName: filepath.Base(fileName),
			}},
	})
	if err != nil {
		return err
	}
	return nil
}

// Data write
func (c *UniClient) DataWrite(ctx context.Context, request *proto.DataWriteRequest) (*proto.EmptyResponse, error) {
	return c.gCli.DataWrite(ctx, request)
}

// Data read
func (c *UniClient) DataRead(ctx context.Context, request *proto.DataReadRequest) (*proto.DataReadResponse, error) {
	return c.gCli.DataRead(ctx, request)
}

// Data print
func (c *UniClient) DataPrint(ctx context.Context, filter *proto.DataReadRequest) {
	dataRes, err := c.DataRead(ctx, filter)
	if err != nil {
		c.logger.Error("Failed to print data", zap.Error(err))
		return
	}

	for _, v := range dataRes.Data {
		switch vv := v.Data.(type) {
		case *proto.DataRecord_CreditCard:
			fmt.Printf("CreditCard:\n\tName: %s\n\tNumber: %s\n\tUntil: %s\n\tHolder: %s",
				vv.CreditCard.Name,
				vv.CreditCard.Number,
				vv.CreditCard.Until,
				vv.CreditCard.Holder)
		case *proto.DataRecord_TextNote:
			fmt.Printf("TextNote:\n\tName: %s\n\tText: %s",
				vv.TextNote.Name,
				vv.TextNote.Text)
		case *proto.DataRecord_Blob:
			fmt.Printf("Files:\n\tName: %s\n\tFileName: %s",
				vv.Blob.Name,
				vv.Blob.FileName)
		case *proto.DataRecord_Credentials:
			fmt.Printf("Credentials:\n\tName: %s\n\tLogin: %s\n\tPassword: %s",
				vv.Credentials.Name,
				vv.Credentials.Login,
				vv.Credentials.Password)
		}

		if len(v.Metadata) > 0 {
			fmt.Print("\n\tMetadata:\n")
		} else {
			fmt.Println("")
		}
		for _, md := range v.Metadata {
			fmt.Printf("\t\t%s: %s\n", md.Name, md.Value)
		}

		fmt.Println("==================")
	}
}
