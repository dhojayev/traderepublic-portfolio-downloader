//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=credentials_interface.go -destination credentials_mock.go -package=auth

package auth

// CredentialsServiceInterface manages authentication credentials.
type CredentialsServiceInterface interface {
	// Load loads credentials from storage
	Load() error

	// Store stores credentials to storage
	Store(token Token) error

	// GetToken returns the current token
	GetToken() Token
}
