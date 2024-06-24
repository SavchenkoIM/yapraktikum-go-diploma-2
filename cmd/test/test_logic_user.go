package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"passwordvault/internal/uni_client"
	"testing"
)

func testLogicUser(ctx context.Context, t *testing.T, client *uni_client.UniClient) {
	var err error

	t.Run("Unregistered_User_Login", func(t *testing.T) {
		_, err = client.UserLogin(ctx, "Victoria", "Victoria's secret")
		assert.Error(t, err)
	})

	t.Run("User_Create", func(t *testing.T) {
		_, err = client.UserCreate(ctx, "Victoria", "Victoria's secret")
		assert.NoError(t, err)
	})

	t.Run("Registered_User_Login_Wrong_Pass", func(t *testing.T) {
		_, err = client.UserLogin(ctx, "Victoria", "Victoria secret")
		assert.Error(t, err)
	})

	t.Run("Registered_User_Login", func(t *testing.T) {
		_, err = client.UserLogin(ctx, "Victoria", "Victoria's secret")
		assert.NoError(t, err)
	})
}
