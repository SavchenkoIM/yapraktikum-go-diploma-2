package main

import (
	"context"
	"go.uber.org/zap"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/uni_client"
)

/*
Will print:

Credentials:
	Name: my_main_creds
	Login: victoria
	Password: victoria's secret
	Metadata:
		site: google.com
==================
Files:
	Name: my_first_script
	FileName: test.py
	Metadata:
		code_quality: bad
==================

And will download 'test.py' to default folder on client PC
*/

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	client := uni_client.NewUniClient(logger)
	client.Start(ctx)

	// Login
	_, err = client.UserCreate(ctx, "victoria", "victoria's secret")
	if err != nil {
		_, err = client.UserLogin(ctx, "victoria", "victoria's secret")
		if err != nil {
			println(err.Error())
		}
	}

	// Write data
	_, err = client.DataWrite(ctx, &proto.DataWriteRequest{
		Action: proto.OperationType_UPSERT,
		Data: &proto.DataWriteRequest_Credentials{Credentials: &proto.DataCredentials{
			Name:     "my_main_creds",
			Login:    "victoria",
			Password: "victoria's secret",
		}},
	})
	if err != nil {
		println(err.Error())
	}

	// Write metadata
	_, err = client.DataWrite(ctx, &proto.DataWriteRequest{
		Action: proto.OperationType_UPSERT,
		Data: &proto.DataWriteRequest_Metadata{Metadata: &proto.MetaDataKV{
			ParentType: proto.DataType_CREDENTIALS,
			ParentName: "my_main_creds",
			Name:       "site",
			Value:      "google.com",
		}},
	})
	if err != nil {
		println(err.Error())
	}

	// Upload file
	err = client.UploadFile(ctx, "my_first_script", "D:\\test.py")
	if err != nil {
		println(err.Error())
	}

	// Write metadata
	_, err = client.DataWrite(ctx, &proto.DataWriteRequest{
		Action: proto.OperationType_UPSERT,
		Data: &proto.DataWriteRequest_Metadata{Metadata: &proto.MetaDataKV{
			ParentType: proto.DataType_BLOB,
			ParentName: "my_first_script",
			Name:       "code_quality",
			Value:      "bad",
		}},
	})
	if err != nil {
		println(err.Error())
	}

	// Print all data
	client.DataPrint(ctx, &proto.DataReadRequest{
		Type:     proto.DataType_UNSPECIFIED,
		NameMask: "%",
		Metadata: nil,
	})

	// Download file
	err = client.DownloadFile(ctx, "my_first_script")
	if err != nil {
		println(err.Error())
	}

	cancel()
}
