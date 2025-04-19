package auth

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
	"github.com/joho/godotenv"
)

const (
	sessionTokenKey = "SESSION_TOKEN"
	refreshTokenKey = "REFRESH_TOKEN"
)

// CredentialsService manages authentication credentials.
type CredentialsService interface {
	// Load loads credentials from storage
	Load() error

	// Store stores credentials to storage
	Store(sessionToken, refreshToken string) error

	// GetSessionToken returns the current session token
	GetSessionToken() string

	// GetRefreshToken returns the current refresh token
	GetRefreshToken() string
}

// FileCredentialsService implements CredentialsService using file storage.
type FileCredentialsService struct {
	filePath     string
	sessionToken string
	refreshToken string
}

// NewFileCredentialsService creates a new file-based credentials service.
// If filePath is empty, it uses the default path from internal.AuthTokenFilename.
func NewFileCredentialsService(filePath string) *FileCredentialsService {
	if filePath == "" {
		filePath = internal.AuthTokenFilename
	}

	return &FileCredentialsService{
		filePath: filePath,
	}
}

// Load loads credentials from a file.
func (s *FileCredentialsService) Load() error {
	env, err := godotenv.Read(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read '%s' file: %w", s.filePath, err)
	}

	s.sessionToken = env[sessionTokenKey]
	s.refreshToken = env[refreshTokenKey]

	return nil
}

// Store stores credentials to a file.
func (s *FileCredentialsService) Store(sessionToken, refreshToken string) error {
	s.sessionToken = sessionToken
	s.refreshToken = refreshToken

	if err := godotenv.Write(map[string]string{
		sessionTokenKey: s.sessionToken,
		refreshTokenKey: s.refreshToken,
	}, s.filePath); err != nil {
		return fmt.Errorf("failed to write tokens to '%s' file: %w", s.filePath, err)
	}

	return nil
}

// GetSessionToken returns the current session token.
func (s *FileCredentialsService) GetSessionToken() string {
	return s.sessionToken
}

// GetRefreshToken returns the current refresh token.
func (s *FileCredentialsService) GetRefreshToken() string {
	return s.refreshToken
}
