package server_store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/storage/file_store"
	"strings"
)

// Handler for data write/delete request
func (s *Storage) DataWrite(ctx context.Context, request *proto.DataWriteRequest) error {

	userId := ctx.Value("LoggedUserId").(string)
	var err error
	var ct pgconn.CommandTag

	switch v := request.Data.(type) {
	case *proto.DataWriteRequest_Credentials:
		switch request.Action {
		case proto.OperationType_UPSERT:
			ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery("credentials"), userId, v.Credentials.Name, v.Credentials.Login, v.Credentials.Password, s.config.Key)
		case proto.OperationType_DELETE:
			ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery("credentials"), userId, v.Credentials.Name)
		}
	case *proto.DataWriteRequest_CreditCard:
		switch request.Action {
		case proto.OperationType_UPSERT:
			ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery("credit_card"), userId, v.CreditCard.Name, v.CreditCard.Number, v.CreditCard.Until, v.CreditCard.Holder, s.config.Key)
		case proto.OperationType_DELETE:
			ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery("credit_card"), userId, v.CreditCard.Name)
		}
	case *proto.DataWriteRequest_TextNote:
		switch request.Action {
		case proto.OperationType_UPSERT:
			ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery("text_note"), userId, v.TextNote.Name, v.TextNote.Text, s.config.Key)
		case proto.OperationType_DELETE:
			ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery("text_note"), userId, v.TextNote.Name)
		}
	case *proto.DataWriteRequest_Blob:
		switch request.Action {
		case proto.OperationType_UPSERT:
			ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery("blob"), userId, v.Blob.Name, v.Blob.FileName, s.config.Key)
		case proto.OperationType_DELETE:
			fileData, err := s.DataRead(ctx, &proto.DataReadRequest{
				Type:     proto.DataType_BLOB,
				NameMask: v.Blob.Name,
			})
			if err != nil {
				return err
			}
			if len(fileData.Data) != 1 {
				return ErrNoDataAffected
			}
			fileName := fileData.Data[0].Data.(*proto.DataRecord_Blob).Blob.FileName
			fskey, err := s.GetFileStoreKey(ctx)
			if err != nil {
				return err
			}
			mc, err := file_store.NewMinioStorage(ctx, s.config.MinioEndPoint, "securestorageservice",
				strings.Replace(userId, "-", "", -1), fskey)
			if err != nil {
				return err
			}
			err = mc.Delete(ctx, fmt.Sprintf(`%s\%s`, strings.Replace(userId, "-", "", -1), fileName))
			if err != nil {
				return err
			}
			ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery("blob"), userId, v.Blob.Name)
		}
	case *proto.DataWriteRequest_Metadata:
		table := ""
		switch v.Metadata.ParentType {
		case proto.DataType_CREDENTIALS:
			table = "credentials"
		case proto.DataType_CREDIT_CARD:
			table = "credit_card"
		case proto.DataType_TEXT_NOTE:
			table = "text_note"
		case proto.DataType_BLOB:
			table = "blob"
		}

		switch request.Action {
		case proto.OperationType_UPSERT:
			ct, err = s.dbConn.Exec(ctx, getMetaDataUpsertQuery(table), userId, v.Metadata.ParentName, v.Metadata.Name, v.Metadata.Value, s.config.Key)
		case proto.OperationType_DELETE:
			ct, err = s.dbConn.Exec(ctx, getMetaDataDeleteQuery(table), userId, v.Metadata.ParentName, v.Metadata.Name)
		}
	default:
		return errors.Wrapf(ErrUnimplemented, "Unknown data type")
	}

	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return ErrNoDataAffected
	}

	return nil
}
