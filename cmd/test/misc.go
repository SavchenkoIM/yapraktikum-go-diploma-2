package test

import (
	"context"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/uni_client"
)

func processCreditCard(
	ctx context.Context,
	client *uni_client.UniClient,
	action proto.OperationType,
	name string,
	number string,
	until string,
	holder string) error {
	_, err := client.DataWrite(ctx, &proto.DataWriteRequest{
		Action: action,
		Data: &proto.DataWriteRequest_CreditCard{CreditCard: &proto.DataCreditCard{
			Name:   name,
			Number: number,
			Until:  until,
			Holder: holder,
		}},
	})
	return err
}

func processTextNote(
	ctx context.Context,
	client *uni_client.UniClient,
	action proto.OperationType,
	name string,
	text string) error {
	_, err := client.DataWrite(ctx, &proto.DataWriteRequest{
		Action: action,
		Data: &proto.DataWriteRequest_TextNote{TextNote: &proto.DataTextNote{
			Name: name,
			Text: text,
		}},
	})
	return err
}

func processMetadata(
	ctx context.Context,
	client *uni_client.UniClient,
	action proto.OperationType,
	pType proto.DataType,
	pName string,
	name string,
	value string) error {
	_, err := client.DataWrite(ctx, &proto.DataWriteRequest{
		Action: action,
		Data: &proto.DataWriteRequest_Metadata{Metadata: &proto.MetaDataKV{
			ParentType: pType,
			ParentName: pName,
			Name:       name,
			Value:      value,
		}},
	})
	return err
}
