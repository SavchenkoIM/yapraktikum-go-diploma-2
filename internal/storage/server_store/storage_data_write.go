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
	switch request.Action {
	case proto.OperationType_UPSERT:
		return s.dataWriteUpsert(ctx, request)
	case proto.OperationType_DELETE:
		return s.dataWriteDelete(ctx, request)
	default:
		return ErrNoDataAffected
	}
}

func (s *Storage) dataWriteDelete(ctx context.Context, request *proto.DataWriteRequest) error {
	var err error
	var ct pgconn.CommandTag

	userId, err := getLoggedUserId(ctx)
	if err != nil {
		return err
	}

	switch v := request.Data.(type) {
	case *proto.DataWriteRequest_Credentials:
		ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery(proto.DataType_CREDENTIALS), userId, v.Credentials.Name)
	case *proto.DataWriteRequest_CreditCard:
		ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery(proto.DataType_CREDIT_CARD), userId, v.CreditCard.Name)
	case *proto.DataWriteRequest_TextNote:
		ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery(proto.DataType_TEXT_NOTE), userId, v.TextNote.Name)
	case *proto.DataWriteRequest_Blob:
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
		ct, err = s.dbConn.Exec(ctx, getDataDeleteQuery(proto.DataType_BLOB), userId, v.Blob.Name)
	case *proto.DataWriteRequest_Metadata:
		ct, err = s.dbConn.Exec(ctx, getMetaDataDeleteQuery(v.Metadata.ParentType), userId, v.Metadata.ParentName, v.Metadata.Name)
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

func (s *Storage) dataWriteUpsert(ctx context.Context, request *proto.DataWriteRequest) error {
	var err error
	var ct pgconn.CommandTag

	userId, err := getLoggedUserId(ctx)
	if err != nil {
		return err
	}

	switch v := request.Data.(type) {
	case *proto.DataWriteRequest_Credentials:
		ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery(proto.DataType_CREDENTIALS), userId, v.Credentials.Name, v.Credentials.Login, v.Credentials.Password, s.config.Key)
	case *proto.DataWriteRequest_CreditCard:
		ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery(proto.DataType_CREDIT_CARD), userId, v.CreditCard.Name, v.CreditCard.Number, v.CreditCard.Until, v.CreditCard.Holder, s.config.Key)
	case *proto.DataWriteRequest_TextNote:
		ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery(proto.DataType_TEXT_NOTE), userId, v.TextNote.Name, v.TextNote.Text, s.config.Key)
	case *proto.DataWriteRequest_Blob:
		ct, err = s.dbConn.Exec(ctx, getDataUpsertQuery(proto.DataType_BLOB), userId, v.Blob.Name, v.Blob.FileName, s.config.Key)
	case *proto.DataWriteRequest_Metadata:
		ct, err = s.dbConn.Exec(ctx, getMetaDataUpsertQuery(v.Metadata.ParentType), userId, v.Metadata.ParentName, v.Metadata.Name, v.Metadata.Value, s.config.Key)
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
