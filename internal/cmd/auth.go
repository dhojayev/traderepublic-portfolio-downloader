package cmd

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
)

func Login(authClient auth.ClientInterface) error {
	resp, err := authClient.Login()
	if err != nil {
		return fmt.Errorf("could not login: %w", err)
	}

	if resp.ProcessID == "" {
		return nil
	}

	var otp string

	fmt.Println("Enter 2FA token:")

	if _, err := fmt.Scanln(&otp); err != nil {
		return fmt.Errorf("could not get otp from input: %w", err)
	}

	if err := authClient.ProvideOTP(resp.ProcessID, otp); err != nil {
		return fmt.Errorf("could not validate otp: %w", err)
	}

	return nil
}
