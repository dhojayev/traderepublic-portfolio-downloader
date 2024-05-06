package console

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
)

type AuthServiceInterface interface {
	Login() error
	SessionToken() api.Token
}

type AuthService struct {
	client      auth.ClientInterface
	phoneNumber string
	pin         string
}

func NewAuthService(client auth.ClientInterface) *AuthService {
	return &AuthService{
		client: client,
	}
}

func (s *AuthService) AcquireCredentials() error {
	var phoneNumber string

	fmt.Print("Enter phone number in international format (+49xxxxxxxxxxxxx): \n")

	if _, err := fmt.Scanln(&phoneNumber); err != nil {
		return fmt.Errorf("could not acquire phone number: %w", err)
	}

	s.phoneNumber = phoneNumber

	pin, err := ReadPassword("pin")
	if err != nil {
		return fmt.Errorf("could not acquire pin: %w", err)
	}

	s.pin = string(pin)

	return nil
}

func (s *AuthService) Login() error {
	if s.phoneNumber == "" || s.pin == "" {
		if err := s.AcquireCredentials(); err != nil {
			return err
		}
	}

	resp, err := s.client.Login(s.phoneNumber, s.pin)
	if err != nil {
		return fmt.Errorf("could not login: %w", err)
	}

	if resp.ProcessID == "" {
		return nil
	}

	input, err := ReadPassword("2FA token")
	if err != nil {
		return fmt.Errorf("could not read otp: %w", err)
	}

	otp := string(input)

	if err := s.client.ProvideOTP(resp.ProcessID, otp); err != nil {
		return fmt.Errorf("could not validate otp: %w", err)
	}

	return nil
}

func (s *AuthService) SessionToken() api.Token {
	return s.client.SessionToken()
}
