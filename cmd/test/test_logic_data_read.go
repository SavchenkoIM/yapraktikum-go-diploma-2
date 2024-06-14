package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/uni_client"
	"slices"
	"testing"
)

func testLogicDataCheck(ctx context.Context, t *testing.T, client *uni_client.UniClient) {

	client.DataPrint(ctx, &proto.DataReadRequest{
		Type:     0,
		NameMask: "%",
		Metadata: nil,
	})

	t.Run("Check_Filter_By_Name_1", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_UNSPECIFIED,
			NameMask: "my%",
			Metadata: nil,
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
			if !slices.Contains([]string{"my_main_credit_card", "my_premium_credit_card", "my_chinese_credit_card"}, name) {
				assert.Fail(t, "unexpected record matched filter")
			}
		}
		assert.Equal(t, 3, len(data.Data))
	})

	t.Run("Check_Filter_By_Name_2", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_UNSPECIFIED,
			NameMask: "%my%",
			Metadata: nil,
		})
		if err != nil {
			assert.Fail(t, err.Error())
		}
		assert.Equal(t, 4, len(data.Data))
	})

	t.Run("Check_Filter_By_Type_1", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_UNSPECIFIED,
			NameMask: "%",
			Metadata: nil,
		})
		assert.NoError(t, err)
		if data != nil {
			assert.Equal(t, 5, len(data.Data))
		}
	})

	t.Run("Check_Filter_By_Type_2", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_TEXT_NOTE,
			NameMask: "%",
			Metadata: nil,
		})
		assert.NoError(t, err)
		if data != nil {
			assert.Equal(t, 1, len(data.Data))
		}
	})

	t.Run("Check_Filter_By_Type_3", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_CREDIT_CARD,
			NameMask: "%",
			Metadata: nil,
		})
		assert.NoError(t, err)
		if data != nil {
			assert.Equal(t, 3, len(data.Data))
		}
	})

	t.Run("Check_Filter_By_Metadata_1", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_CREDIT_CARD,
			NameMask: "%",
			Metadata: []*proto.MetaDataKV{
				{Name: "type", Value: "%M%C%"},
			},
		})
		assert.NoError(t, err)
		if data != nil {
			assert.Equal(t, 1, len(data.Data))
			cc, ok := data.Data[0].Data.(*proto.DataRecord_CreditCard)
			assert.True(t, ok)
			if ok {
				assert.Equal(t, "my_premium_credit_card", cc.CreditCard.Name)
			}
		}
	})

	t.Run("Check_Filter_By_Metadata_2", func(t *testing.T) {
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
			if !slices.Contains([]string{"my_premium_credit_card", "my_chinese_credit_card"}, name) {
				assert.Fail(t, "unexpected record matched filter")
			}
		}
	})

	t.Run("Check_Filter_By_Metadata_And_Name", func(t *testing.T) {
		data, err := client.DataRead(ctx, &proto.DataReadRequest{
			Type:     proto.DataType_CREDIT_CARD,
			NameMask: "%chinese%",
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
			if !slices.Contains([]string{"my_chinese_credit_card"}, name) {
				assert.Fail(t, "unexpected record matched filter")
			}
		}
	})
}
