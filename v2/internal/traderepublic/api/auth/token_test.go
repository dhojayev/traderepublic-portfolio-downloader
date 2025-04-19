package auth_test

import (
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	t.Parallel()
	t.Run("NewToken", func(t *testing.T) {
		t.Parallel()

		token := auth.NewToken()

		assert.Empty(t, token.SessionToken(), "Session token should be empty")
		assert.Empty(t, token.RefreshToken(), "Refresh token should be empty")
	})

	t.Run("NewTokenWithValues", func(t *testing.T) {
		t.Parallel()

		sessionToken := "test-session-token"
		refreshToken := "test-refresh-token"

		token := auth.NewTokenWithValues(sessionToken, refreshToken)

		assert.Equal(t, sessionToken, token.SessionToken(), "Session token should match")
		assert.Equal(t, refreshToken, token.RefreshToken(), "Refresh token should match")
	})

	t.Run("EmptyValues", func(t *testing.T) {
		t.Parallel()

		token := auth.NewTokenWithValues("", "")

		assert.Empty(t, token.SessionToken(), "Session token should be empty")
		assert.Empty(t, token.RefreshToken(), "Refresh token should be empty")
	})
}
