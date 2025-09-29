package main

import (
	"fmt"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
)

type App struct {
	authClient         *auth.Client
	credentialsService auth.CredentialsServiceInterface
	log                *slog.Logger
}

func NewApp(
	authClient *auth.Client,
	credentialsService auth.CredentialsServiceInterface,
	log *slog.Logger,
) App {
	return App{
		authClient:         authClient,
		credentialsService: credentialsService,
		log:                log,
	}
}

func (a *App) Run() error {
	err := a.credentialsService.Load()
	if err != nil {
		a.log.Warn("Failed to load credentials, need to authenticate", "error", err)

		err := a.authenticate()
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	a.log.Info("Credentials loaded successfully")

	return nil
}

func (a *App) authenticate() error {
	token, err := a.authClient.Login()
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	err = a.credentialsService.Store(token)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	a.log.Info("Authentication successful")

	return nil
}
