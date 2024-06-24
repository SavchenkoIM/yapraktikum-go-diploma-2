package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/uni_client"
	"slices"
	"testing"
)

func testLogicDataDelete(ctx context.Context, t *testing.T, client *uni_client.UniClient) {

	t.Run("Delete_Non_Existing_File", func(t *testing.T) {
		err := client.DeleteFile(ctx, "test__file")
		assert.Error(t, err)
	})

	t.Run("Delete_Existing_File", func(t *testing.T) {
		err := client.DeleteFile(ctx, "test_file")
		assert.NoError(t, err)
	})

	t.Run("Check_No_Files_Remaining", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_BLOB,
			NameMask: "%",
			Metadata: nil,
		})
		if err != nil {
			assert.Fail(t, err.Error())
		}
		assert.Equal(t, 0, len(data.Data))
	})

	t.Run("Delete_Metadata_Wrong_Parent_Type", func(t *testing.T) {
		err := processMetadata(ctx, client, proto.OperationType_DELETE, proto.DataType_TEXT_NOTE,
			"my_chinese_credit_card", "type", "")
		assert.Error(t, err)
	})

	t.Run("Delete_Metadata", func(t *testing.T) {
		err := processMetadata(ctx, client, proto.OperationType_DELETE, proto.DataType_CREDIT_CARD,
			"my_chinese_credit_card", "type", "")
		assert.NoError(t, err)
	})

	t.Run("ReCheck_Filter_By_Metadata", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_CREDIT_CARD,
			NameMask: "%",
			Metadata: []*proto.MetaDataKV{
				{Name: "type", Value: "%M%C%"},
				{Name: "type", Value: "Union%Pay"},
			},
		})
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		var name string
		for _, data := range data.Data {
			cc := data.GetCreditCard()
			if cc != nil {
				name = cc.Name
			} else {
				assert.Fail(t, "unexpected record matched filter")
			}
			if !slices.Contains([]string{"my_premium_credit_card"}, name) {
				assert.Fail(t, "unexpected record matched filter")
			}
		}
	})

	t.Run("Delete_CreditCard_Record", func(t *testing.T) {
		err := processCreditCard(ctx, client, proto.OperationType_DELETE, "my_premium_credit_card", "", "", "")
		assert.NoError(t, err)
	})

	t.Run("ReCheck_Filter_By_Name", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_CREDIT_CARD,
			NameMask: "my_premium_credit_card",
			Metadata: nil,
		})
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, 0, len(data.Data))
	})
}
