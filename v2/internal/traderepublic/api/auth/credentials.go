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

// CredentialsService manages authentication credentials
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

// FileCredentialsService implements CredentialsService using file storage
type FileCredentialsService struct {
	sessionToken string
	refreshToken string
}

// NewFileCredentialsService creates a new file-based credentials service
func NewFileCredentialsService() *FileCredentialsService {
	return &FileCredentialsService{}
}

// Load loads credentials from a file
func (s *FileCredentialsService) Load() error {
	env, err := godotenv.Read(internal.AuthTokenFilename)
	if err != nil {
		return fmt.Errorf("failed to read '%s' file: %w", internal.AuthTokenFilename, err)
	}

	s.sessionToken = env[sessionTokenKey]
	s.refreshToken = env[refreshTokenKey]

	return nil
}

// Store stores credentials to a file
func (s *FileCredentialsService) Store(sessionToken, refreshToken string) error {
	s.sessionToken = sessionToken
	s.refreshToken = refreshToken

	if err := godotenv.Write(map[string]string{
		sessionTokenKey: s.sessionToken,
		refreshTokenKey: s.refreshToken,
	}, internal.AuthTokenFilename); err != nil {
		return fmt.Errorf("failed to write tokens to '%s' file: %w", internal.AuthTokenFilename, err)
	}

	return nil
}

// GetSessionToken returns the current session token
func (s *FileCredentialsService) GetSessionToken() string {
	return s.sessionToken
}

// GetRefreshToken returns the current refresh token
func (s *FileCredentialsService) GetRefreshToken() string {
	return s.refreshToken
}
