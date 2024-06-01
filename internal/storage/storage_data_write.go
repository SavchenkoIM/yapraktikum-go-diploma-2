package storage

import (
	"context"
	proto "passwordvault/internal/proto/gen"
)

func (s *Storage) DataWrite(ctx context.Context, request *proto.DataWriteRequest) error {

	userId := ctx.Value("LoggedUserId").(string)
	var err error

	switch v := request.Data.(type) {
	case *proto.DataWriteRequest_Credentials:
		switch request.Action {
		case proto.OperationType_UPSERT:
			_, err = s.dbConn.Exec(ctx, getDataUpsertQuery("credentials"), userId, v.Credentials.Name, v.Credentials.Login, v.Credentials.Password, s.config.Key)
		case proto.OperationType_DELETE:
			_, err = s.dbConn.Exec(ctx, getDataDeleteQuery("credentials"), userId, v.Credentials.Name)
		}
	case *proto.DataWriteRequest_CreditCard:
		switch request.Action {
		case proto.OperationType_UPSERT:
			_, err = s.dbConn.Exec(ctx, getDataUpsertQuery("credit_card"), userId, v.CreditCard.Name, v.CreditCard.Number, v.CreditCard.Until, v.CreditCard.Holder, s.config.Key)
		case proto.OperationType_DELETE:
			_, err = s.dbConn.Exec(ctx, getDataDeleteQuery("credit_card"), userId, v.CreditCard.Name)
		}
	case *proto.DataWriteRequest_TextNote:
		switch request.Action {
		case proto.OperationType_UPSERT:
			_, err = s.dbConn.Exec(ctx, getDataUpsertQuery("text_note"), userId, v.TextNote.Name, v.TextNote.Text, s.config.Key)
		case proto.OperationType_DELETE:
			_, err = s.dbConn.Exec(ctx, getDataDeleteQuery("text_note"), userId, v.TextNote.Name)
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
		}

		switch request.Action {
		case proto.OperationType_UPSERT:
			_, err = s.dbConn.Exec(ctx, getMetaDataUpsertQuery(table), userId, v.Metadata.ParentName, v.Metadata.Name, v.Metadata.Value, s.config.Key)
		case proto.OperationType_DELETE:
			_, err = s.dbConn.Exec(ctx, getMetaDataDeleteQuery(table), userId, v.Metadata.ParentName, v.Metadata.Name)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
