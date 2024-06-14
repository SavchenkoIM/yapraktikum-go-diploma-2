package server_store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "passwordvault/internal/proto/gen"
)

// Handler for data read request
func (s *Storage) DataRead(ctx context.Context, request *proto.DataReadRequest) (*proto.DataReadResponse, error) {

	res := &proto.DataReadResponse{}

	LoggedUserId := ctx.Value("LoggedUserId").(string)

	var tableName string

	var dataTypes []proto.DataType
	if request.Type != proto.DataType_UNSPECIFIED {
		dataTypes = []proto.DataType{request.Type}
	} else {
		dataTypes = []proto.DataType{
			proto.DataType_CREDENTIALS,
			proto.DataType_CREDIT_CARD,
			proto.DataType_TEXT_NOTE,
			proto.DataType_BLOB,
		}
	}

	for _, dataType := range dataTypes {

		switch dataType {
		case proto.DataType_CREDENTIALS:
			tableName = "credentials"
		case proto.DataType_CREDIT_CARD:
			tableName = "credit_card"
		case proto.DataType_TEXT_NOTE:
			tableName = "text_note"
		case proto.DataType_BLOB:
			tableName = "blob"
		default:
			return nil, errors.Wrapf(ErrUnimplemented, "Unknown data type")
		}

		query, params := getDataReadQuery(tableName, LoggedUserId, request, s.config)

		rows, err := s.dbConn.Query(ctx, query, params...)
		if err != nil {
			return nil, err
		}

		var (
			uuid pgtype.UUID
			name string
		)

		for rows.Next() {
			dr := &proto.DataRecord{}
			switch dataType {
			case proto.DataType_CREDENTIALS:
				var (
					login string
					pass  string
				)
				err := rows.Scan(&uuid, &name, &login, &pass)
				if err != nil {
					return nil, status.Error(codes.Unknown, fmt.Sprintf("database error: %v", err))
				}
				dr = &proto.DataRecord{
					Data: &proto.DataRecord_Credentials{Credentials: &proto.DataCredentials{
						Name:     name,
						Login:    login,
						Password: pass,
					}},
					Metadata: make([]*proto.MetaDataKV, 0),
				}
			case proto.DataType_CREDIT_CARD:
				var (
					number string
					until  string
					holder string
				)
				err := rows.Scan(&uuid, &name, &number, &until, &holder)
				if err != nil {
					return nil, status.Error(codes.Unknown, fmt.Sprintf("database error: %v", err))
				}
				dr = &proto.DataRecord{
					Data: &proto.DataRecord_CreditCard{CreditCard: &proto.DataCreditCard{
						Name:   name,
						Number: until,
						Until:  until,
						Holder: holder,
					}},
					Metadata: make([]*proto.MetaDataKV, 0),
				}
			case proto.DataType_TEXT_NOTE:
				var (
					text string
				)
				err := rows.Scan(&uuid, &name, &text)
				if err != nil {
					return nil, status.Error(codes.Unknown, fmt.Sprintf("database error: %v", err))
				}
				dr = &proto.DataRecord{
					Data: &proto.DataRecord_TextNote{TextNote: &proto.DataTextNote{
						Name: name,
						Text: text,
					}},
					Metadata: make([]*proto.MetaDataKV, 0),
				}
			case proto.DataType_BLOB:
				var (
					fileName string
				)
				err := rows.Scan(&uuid, &name, &fileName)
				if err != nil {
					return nil, status.Error(codes.Unknown, fmt.Sprintf("database error: %v", err))
				}
				dr = &proto.DataRecord{
					Data: &proto.DataRecord_Blob{Blob: &proto.DataBLOB{
						Name:     name,
						FileName: fileName,
					}},
					Metadata: make([]*proto.MetaDataKV, 0),
				}
			default:
				continue
			}

			query, params := getMetadataReadQuery(tableName, uuid, s.config)
			subRows, err := s.dbConn.Query(ctx, query, params...)
			if err != nil {
				return nil, err
			}

			var (
				mName  string
				mValue string
			)
			for subRows.Next() {
				err = subRows.Scan(&mName, &mValue)
				if err != nil {
					return nil, status.Error(codes.Unknown, fmt.Sprintf("database error: %v", err))
				}
				dr.Metadata = append(dr.Metadata, &proto.MetaDataKV{
					ParentType: dataType,
					ParentName: name,
					Name:       mName,
					Value:      mValue,
				})
			}

			res.Data = append(res.Data, dr)
		}

	}

	return res, nil

}
