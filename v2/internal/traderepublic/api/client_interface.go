//go:generate go tool mockgen -source=client_interface.go -destination client_mock.go -package=api

package api

import (
	"net/http"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

// ClientInterface is the interface for the Trade Republic API client.
type ClientInterface interface {
	// Login logs in with phone number and PIN.
	Login(requestBody traderepublic.APILoginRequest) (string, error)

	// PostOTP verifies the OTP.
	PostOTP(processID, otp string) ([]*http.Cookie, error)
}
