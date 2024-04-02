package cmd

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/util"
)

func Login(authClient *auth.Client) error {
	resp, err := authClient.Login()
	if err != nil {
		return fmt.Errorf("could not login: %w", err)
	}

	if resp.ProcessID == "" {
		return nil
	}

	input, err := util.ReadPassword("2FA token")
	if err != nil {
		return fmt.Errorf("could not read otp: %w", err)
	}

	otp := string(input)

	if err := authClient.ProvideOTP(resp.ProcessID, otp); err != nil {
		return fmt.Errorf("could not validate otp: %w", err)
	}

	return nil
}
