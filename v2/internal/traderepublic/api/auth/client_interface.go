//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client_interface.go -destination client_mock.go -package=auth

package auth

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/restclient"
)

type ClientInterface interface {
	Login(phoneNumber, pin string) (restclient.APILoginResponse, error)
	ProvideOTP(processID, otp string) error
	SessionToken() api.Token
}
