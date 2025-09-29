package message

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message/subscriber"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
)

const (
	// Subscription types.
	TypeTimelineTransactions = "timelineTransactions"
	TypeTimelineDetail       = "timelineDetailV2"
)

type ClientInterface interface {
	SubscribeToTimelineTransactions(ctx context.Context) error
}

type Client struct {
	credentialsService auth.CredentialsServiceInterface
	wsClient           websocketclient.ClientInterface
	logger             *slog.Logger
}

func NewClient(credentialsService auth.CredentialsServiceInterface, wsClient websocketclient.ClientInterface, logger *slog.Logger) *Client {
	return &Client{
		credentialsService: credentialsService,
		wsClient:           wsClient,
		logger:             logger,
	}
}

// SubscribeToTimelineTransactions subscribes to timeline transactions data.
func (c *Client) SubscribeToTimelineTransactions(ctx context.Context) error {
	ch, err := c.SubscribeToTimelineTransactionsWithCursor(ctx, "")
	if err != nil {
		return err
	}

	var response struct {
		Cursors struct {
			After *string `json:"after",omitempty`
		} `json:"cursors"`
	}

	unifiedChannel := make(chan []byte, 1)
	data := <-ch
	unifiedChannel <- data

	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	go func() {
		for response.Cursors.After != nil {
			ch, err = c.SubscribeToTimelineTransactionsWithCursor(ctx, *response.Cursors.After)
			if err != nil {
				c.logger.Error("error subscribing to timeline transactions", "error", err)

				return
			}

			data = <-ch

			err = json.Unmarshal(data, &response)
			if err != nil {
				c.logger.Error("error subscribing to timeline transactions", "error", err)

				return
			}

			unifiedChannel <- data
		}
	}()

	subscriber := subscriber.NewSubscriber("timelineTransactions", unifiedChannel, writer.NewResponseWriter(), c.logger)
	subscriber.Listen()

	return nil
}

// SubscribeToTimelineTransactionsWithCursor subscribes to timeline transactions data with a cursor.
func (c *Client) SubscribeToTimelineTransactionsWithCursor(ctx context.Context, cursor string) (<-chan []byte, error) {
	params := map[string]any{}

	// Add cursor if provided
	if cursor != "" {
		params["after"] = cursor
	}

	data := c.prepareSubscription(TypeTimelineTransactions, params)

	c.wsClient.Connect(ctx)

	return c.wsClient.Subscribe(ctx, data)
}

// SubscribeToTimelineDetail subscribes to timeline detail data.
func (c *Client) SubscribeToTimelineDetail(ctx context.Context, itemID string) (<-chan []byte, error) {
	data := c.prepareSubscription(TypeTimelineDetail, map[string]any{"id": itemID})

	c.wsClient.Connect(ctx)

	return c.wsClient.Subscribe(ctx, data)
}

// prepareSubscription prepares a subscription request with the given parameters.
func (c *Client) prepareSubscription(dataType string, params map[string]any) map[string]any {
	data := map[string]any{
		"type":  dataType,
		"token": c.credentialsService.GetToken().Session(),
	}

	// Add additional parameters
	for k, v := range params {
		data[k] = v
	}

	return data
}
