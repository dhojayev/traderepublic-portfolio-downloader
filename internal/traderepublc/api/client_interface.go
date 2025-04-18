//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client_interface.go -destination client_mock.go -package=api

package api

// ClientInterface is the interface for the Trade Republic API client.
type ClientInterface interface {
	// Login logs in with phone number and PIN.
	Login(requestBody LoginRequest, refreshToken Token) (LoginResponse, Token, error)

	// PostOTP verifies the OTP.
	PostOTP(processID, otp string) (Token, Token, error)

	// Session refreshes the session token.
	Session(refreshToken Token) (Token, error)
}
