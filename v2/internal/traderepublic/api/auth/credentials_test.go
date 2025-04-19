package auth_test

import (
	"path/filepath"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileCredentialsService(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Each test will use its own credentials file path

	// Test data
	testSessionToken := "test-session-token"
	testRefreshToken := "test-refresh-token"

	t.Run("Store and Load", func(t *testing.T) {
		t.Parallel()
		// Create a unique file path for this test
		testFile := filepath.Join(tempDir, "store-and-load.env")

		// Create a new credentials service with the test file path
		service := auth.NewFileCredentialsService(testFile)

		// Store the credentials
		err := service.Store(testSessionToken, testRefreshToken)
		require.NoError(t, err, "Storing credentials should not error")

		// Create a new service to load the credentials
		loadService := auth.NewFileCredentialsService(testFile)

		// Load the credentials
		err = loadService.Load()
		require.NoError(t, err, "Loading credentials should not error")

		// Check that the loaded credentials match the stored ones
		assert.Equal(t, testSessionToken, loadService.GetSessionToken(), "Session token should match")
		assert.Equal(t, testRefreshToken, loadService.GetRefreshToken(), "Refresh token should match")
	})

	t.Run("Load Non-existent File", func(t *testing.T) {
		t.Parallel()
		// Create a path to a non-existent file
		nonExistentFile := filepath.Join(tempDir, "non-existent.env")

		// Create a new service with the non-existent file path
		service := auth.NewFileCredentialsService(nonExistentFile)

		// Try to load the credentials
		err := service.Load()
		assert.Error(t, err, "Loading non-existent file should return an error")
	})

	t.Run("Store Empty Tokens", func(t *testing.T) {
		t.Parallel()
		// Create a unique file path for this test
		testFile := filepath.Join(tempDir, "store-empty.env")

		// Create a new credentials service with the test file path
		service := auth.NewFileCredentialsService(testFile)

		// Store empty tokens
		err := service.Store("", "")
		require.NoError(t, err, "Storing empty tokens should not error")

		// Create a new service to load the credentials
		loadService := auth.NewFileCredentialsService(testFile)

		// Load the credentials
		err = loadService.Load()
		require.NoError(t, err, "Loading credentials should not error")

		// Check that the loaded credentials are empty
		assert.Empty(t, loadService.GetSessionToken(), "Session token should be empty")
		assert.Empty(t, loadService.GetRefreshToken(), "Refresh token should be empty")
	})

	t.Run("Overwrite Existing Tokens", func(t *testing.T) {
		t.Parallel()
		// Create a unique file path for this test
		testFile := filepath.Join(tempDir, "overwrite.env")

		// Create a new credentials service with the test file path
		service := auth.NewFileCredentialsService(testFile)

		// Store initial tokens
		err := service.Store(testSessionToken, testRefreshToken)
		require.NoError(t, err, "Storing initial tokens should not error")

		// Store new tokens
		newSessionToken := "new-session-token"
		newRefreshToken := "new-refresh-token"

		err = service.Store(newSessionToken, newRefreshToken)
		require.NoError(t, err, "Storing new tokens should not error")

		// Create a new service to load the credentials
		loadService := auth.NewFileCredentialsService(testFile)

		// Load the credentials
		err = loadService.Load()
		require.NoError(t, err, "Loading credentials should not error")

		// Check that the loaded credentials match the new tokens
		assert.Equal(t, newSessionToken, loadService.GetSessionToken(), "Session token should match new token")
		assert.Equal(t, newRefreshToken, loadService.GetRefreshToken(), "Refresh token should match new token")
	})
}
