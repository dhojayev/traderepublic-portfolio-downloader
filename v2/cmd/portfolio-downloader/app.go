package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message"
)

type App struct {
	authClient         *auth.Client
	credentialsService auth.CredentialsServiceInterface
	messageClient      message.ClientInterface
	eventBus           *bus.EventBus
}

func NewApp(
	authClient *auth.Client,
	credentialsService auth.CredentialsServiceInterface,
	messageClient message.ClientInterface,
	eventBus *bus.EventBus,
) App {
	return App{
		authClient:         authClient,
		credentialsService: credentialsService,
		messageClient:      messageClient,
		eventBus:           eventBus,
	}
}

func (a *App) Run() error {
	err := a.credentialsService.Load()
	if err != nil {
		slog.Warn("Failed to load credentials, need to authenticate", "error", err)

		err := a.authenticate()
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	slog.Info("Starting downloading transactions")

	err = a.messageClient.SubscribeToTimelineTransactions(context.Background())
	if err != nil {
		return fmt.Errorf("subscription failed: %w", err)
	}

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

	slog.Info("Authentication successful")

	return nil
}
