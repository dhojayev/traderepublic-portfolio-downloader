package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message/subscriber"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

type App struct {
	authClient         *auth.Client
	credentialsService auth.CredentialsServiceInterface
	messageClient      message.ClientInterface
}

func NewApp(
	authClient *auth.Client,
	credentialsService auth.CredentialsServiceInterface,
	messageClient message.ClientInterface,
) App {
	return App{
		authClient:         authClient,
		credentialsService: credentialsService,
		messageClient:      messageClient,
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

	ch, err := a.messageClient.SubscribeToTimelineTransactions(context.Background())
	if err != nil {
		return fmt.Errorf("subscription failed: %w", err)
	}

	sub := subscriber.NewSubscriber("timelineTransactions", ch, writer.NewResponseWriter())
	sub.Listen()

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
