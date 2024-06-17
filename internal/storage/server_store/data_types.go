package server_store

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	proto "passwordvault/internal/proto/gen"
)

type objectTypeDescription struct {
	FieldsCount int
	TableName   string
	ScanFunc    func(pgx.Rows, *string, *pgtype.Text, *pgtype.Text) (*proto.DataRecord, error)
}

var supportedObjectTypes = []proto.DataType{
	proto.DataType_CREDENTIALS,
	proto.DataType_CREDIT_CARD,
	proto.DataType_TEXT_NOTE,
	proto.DataType_BLOB}

var objectTypes = map[proto.DataType]objectTypeDescription{
	proto.DataType_CREDENTIALS: {
		FieldsCount: 2,
		TableName:   "credentials",
		ScanFunc: func(rows pgx.Rows, n *string, mn *pgtype.Text, mc *pgtype.Text) (*proto.DataRecord, error) {
			var c []string = make([]string, 10)
			err := rows.Scan(n, &c[0], &c[1], mn, mc)
			if err != nil {
				return &proto.DataRecord{}, err
			}
			return &proto.DataRecord{
				Data: &proto.DataRecord_Credentials{Credentials: &proto.DataCredentials{
					Name:     *n,
					Login:    c[0],
					Password: c[1],
				}},
			}, nil
		},
	},
	proto.DataType_CREDIT_CARD: {
		FieldsCount: 3,
		TableName:   "credit_card",
		ScanFunc: func(rows pgx.Rows, n *string, mn *pgtype.Text, mc *pgtype.Text) (*proto.DataRecord, error) {
			var c []string = make([]string, 10)
			err := rows.Scan(n, &c[0], &c[1], &c[2], mn, mc)
			if err != nil {
				return &proto.DataRecord{}, err
			}
			return &proto.DataRecord{
				Data: &proto.DataRecord_CreditCard{CreditCard: &proto.DataCreditCard{
					Name:   *n,
					Number: c[0],
					Until:  c[1],
					Holder: c[2],
				}},
			}, nil
		},
	},
	proto.DataType_TEXT_NOTE: {
		FieldsCount: 1,
		TableName:   "text_note",
		ScanFunc: func(rows pgx.Rows, n *string, mn *pgtype.Text, mc *pgtype.Text) (*proto.DataRecord, error) {
			var c []string = make([]string, 10)
			err := rows.Scan(n, &c[0], mn, mc)
			if err != nil {
				return &proto.DataRecord{}, err
			}
			return &proto.DataRecord{
				Data: &proto.DataRecord_TextNote{TextNote: &proto.DataTextNote{
					Name: *n,
					Text: c[0],
				}},
			}, nil
		},
	},
	proto.DataType_BLOB: {
		FieldsCount: 1,
		TableName:   "blob",
		ScanFunc: func(rows pgx.Rows, n *string, mn *pgtype.Text, mc *pgtype.Text) (*proto.DataRecord, error) {
			var c []string = make([]string, 10)
			err := rows.Scan(n, &c[0], mn, mc)
			if err != nil {
				return &proto.DataRecord{}, err
			}
			return &proto.DataRecord{
				Data: &proto.DataRecord_Blob{Blob: &proto.DataBLOB{
					Name:     *n,
					FileName: c[0],
				}},
			}, nil
		},
	},
}
