//go:generate go tool mockgen -source=client_interface.go -destination client_mock.go -package=auth

package auth

type (
	PhoneNumber string
	Pin         string
	OTP         string
	ProcessID   string
)

type ClientInterface interface {
	Login(phoneNumber PhoneNumber, pin Pin) (ProcessID, error)
	ProvideOTP(processID ProcessID, otp OTP) (Token, error)
}
