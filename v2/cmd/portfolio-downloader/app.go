package main

import (
	"fmt"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
)

type App struct {
	inputHandler console.InputHandlerInterface
	authClient   *auth.Client
	log          *slog.Logger
}

func NewApp(
	inputHandler console.InputHandlerInterface,
	authClient *auth.Client,
	log *slog.Logger,
) App {
	return App{
		inputHandler: inputHandler,
		authClient:   authClient,
		log:          log,
	}
}

func (a *App) Run() error {
	if err := a.authenticate(); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	a.log.Info("Authentication successful")
	return nil
}

func (a *App) authenticate() error {
	phoneNumber, err := a.inputHandler.GetPhoneNumber()
	if err != nil {
		return fmt.Errorf("failed to get phone number: %w", err)
	}

	pin, err := a.inputHandler.GetPIN()
	if err != nil {
		return fmt.Errorf("failed to get PIN: %w", err)
	}

	processID, err := a.authClient.Login(auth.PhoneNumber(phoneNumber), auth.Pin(pin))
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	otp, err := a.inputHandler.GetOTP()
	if err != nil {
		return fmt.Errorf("failed to get OTP: %w", err)
	}

	token, err := a.authClient.ProvideOTP(processID, auth.OTP(otp))
	if err != nil {
		return fmt.Errorf("OTP validation failed: %w", err)
	}

	a.log.Debug("Obtained tokens", "session", token.Session(), "refresh", token.Refresh())
	// Here you can store the token securely for future use
	// For example, save it to a file or a secure storage

	return err
}
