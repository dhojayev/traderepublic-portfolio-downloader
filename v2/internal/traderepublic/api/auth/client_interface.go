//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client_interface.go -destination client_mock.go -package=auth

package auth

type (
	PhoneNumber string
	Pin         string
	OTP         string
	ProcessID   string
)

type ClientInterface interface {
	Login(PhoneNumber, Pin) (ProcessID, error)
	ProvideOTP(ProcessID, OTP) (Token, error)
}
