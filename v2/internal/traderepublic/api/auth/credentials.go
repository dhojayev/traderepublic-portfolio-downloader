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

// FileCredentialsService implements CredentialsServiceInterface using file storage.
type FileCredentialsService struct {
	filePath string
	token    Token
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

	s.token = NewTokenWithValues(env[sessionTokenKey], env[refreshTokenKey])

	return nil
}

// Store stores credentials to a file.
func (s *FileCredentialsService) Store(token Token) error {
	s.token = token

	if err := godotenv.Write(map[string]string{
		sessionTokenKey: token.Session(),
		refreshTokenKey: token.Refresh(),
	}, s.filePath); err != nil {
		return fmt.Errorf("failed to write tokens to '%s' file: %w", s.filePath, err)
	}

	return nil
}

// GetToken returns the current token.
func (s *FileCredentialsService) GetToken() Token {
	return s.token
}
