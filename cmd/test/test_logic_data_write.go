package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/uni_client"
	"testing"
)

func testLogicDataWrite(ctx context.Context, t *testing.T, client *uni_client.UniClient) {
	var err error

	t.Run("Add_Data_Wrong_Action", func(t *testing.T) {
		err = processCreditCard(ctx, client, 0, "my_main_credit_card",
			"1234567890123456", "1214", "Victoria Vysotskaya")
		assert.Error(t, err)
	})

	t.Run("Add_Credit_Card", func(t *testing.T) {
		err = processCreditCard(ctx, client, proto.OperationType_UPSERT, "my_main_credit_card",
			"1234567890123456", "1214", "Victoria Vysotskaya")
		assert.NoError(t, err)
	})

	t.Run("Add_Metadata_For_Wrong_Record", func(t *testing.T) {
		err = processMetadata(ctx, client, proto.OperationType_UPSERT, proto.DataType_TEXT_NOTE,
			"my_main_credit_card", "type", "VISA")
		assert.Error(t, err)
	})

	t.Run("Add_Metadata_For_Right_Record", func(t *testing.T) {
		err = processMetadata(ctx, client, proto.OperationType_UPSERT, proto.DataType_CREDIT_CARD,
			"my_main_credit_card", "type", "VISA")
		assert.NoError(t, err)
	})

	t.Run("Add_More_Metadata_For_Right_Record", func(t *testing.T) {
		err = processMetadata(ctx, client, proto.OperationType_UPSERT, proto.DataType_CREDIT_CARD,
			"my_main_credit_card", "class", "GOLD")
		assert.NoError(t, err)
	})

	t.Run("Add_Second_Credit_Card", func(t *testing.T) {
		err = processCreditCard(ctx, client, proto.OperationType_UPSERT, "my_premium_credit_card",
			"0987654321098765", "1214", "Victoria Vysotskaya")
		assert.NoError(t, err)
	})

	t.Run("Add_Metadata_For_Second_Record", func(t *testing.T) {
		err = processMetadata(ctx, client, proto.OperationType_UPSERT, proto.DataType_CREDIT_CARD,
			"my_premium_credit_card", "type", "MasterCard")
		assert.NoError(t, err)
	})

	t.Run("Add_More_Metadata_For_Second_Record", func(t *testing.T) {
		err = processMetadata(ctx, client, proto.OperationType_UPSERT, proto.DataType_CREDIT_CARD,
			"my_premium_credit_card", "class", "WORLD ELITE")
		assert.NoError(t, err)
	})

	t.Run("Add_Exotic_Credit_Card", func(t *testing.T) {
		err = processCreditCard(ctx, client, proto.OperationType_UPSERT, "my_chinese_credit_card",
			"7890992665382845", "1214", "Victoria Vysotskaya")
		assert.NoError(t, err)
	})

	t.Run("Add_Metadata_For_Exotic_Record", func(t *testing.T) {
		err = processMetadata(ctx, client, proto.OperationType_UPSERT, proto.DataType_CREDIT_CARD,
			"my_chinese_credit_card", "type", "UnionPay")
		assert.NoError(t, err)
	})

	t.Run("Add_Text_Note_Record", func(t *testing.T) {
		err = processTextNote(ctx, client, proto.OperationType_UPSERT, "not_my_text_note", "Hello, world!")
		assert.NoError(t, err)
	})

}
