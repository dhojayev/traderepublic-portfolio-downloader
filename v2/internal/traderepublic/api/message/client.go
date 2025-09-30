package message

import (
	"context"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type ClientInterface interface {
	SubscribeToTimelineTransactions(ctx context.Context) (<-chan []byte, error)
	SubscribeToTimelineDetail(ctx context.Context, itemID string) (<-chan []byte, error)
}

type Client struct {
	credentialsService auth.CredentialsServiceInterface
	wsClient           websocketclient.ClientInterface
}

func NewClient(credentialsService auth.CredentialsServiceInterface, wsClient websocketclient.ClientInterface) *Client {
	return &Client{
		credentialsService: credentialsService,
		wsClient:           wsClient,
	}
}

// SubscribeToTimelineTransactions subscribes to timeline transactions data.
func (c *Client) SubscribeToTimelineTransactions(ctx context.Context) (<-chan []byte, error) {
	ch, err := c.SubscribeToTimelineTransactionsWithCursor(ctx, "")
	if err != nil {
		return nil, err
	}

	unifiedChannel := make(chan []byte, 1)
	data := <-ch
	unifiedChannel <- data

	var response traderepublic.TimelineTransactionsSchemaJson

	err = response.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	go func() {
		for response.Cursors.After != nil {
			ch, err = c.SubscribeToTimelineTransactionsWithCursor(ctx, *response.Cursors.After)
			if err != nil {
				slog.Error("error subscribing to timeline transactions", "error", err)

				return
			}

			data = <-ch

			err = response.UnmarshalJSON(data)
			if err != nil {
				slog.Error("error subscribing to timeline transactions", "error", err)

				return
			}

			unifiedChannel <- data
		}
	}()

	return unifiedChannel, nil
}

// SubscribeToTimelineTransactionsWithCursor subscribes to timeline transactions data with a cursor.
func (c *Client) SubscribeToTimelineTransactionsWithCursor(ctx context.Context, cursor string) (<-chan []byte, error) {
	data := traderepublic.WebsocketSubRequestSchemaJson{
		Token: c.credentialsService.GetToken().Session(),
		Type:  traderepublic.WebsocketSubRequestSchemaJsonTypeTimelineTransactions,
		After: &cursor,
	}

	c.wsClient.Connect(ctx)

	return c.wsClient.Subscribe(ctx, data)
}

// SubscribeToTimelineDetail subscribes to timeline detail data.
func (c *Client) SubscribeToTimelineDetail(ctx context.Context, itemID string) (<-chan []byte, error) {
	data := traderepublic.WebsocketSubRequestSchemaJson{
		Id:    &itemID,
		Token: c.credentialsService.GetToken().Session(),
		Type:  traderepublic.WebsocketSubRequestSchemaJsonTypeTimelineDetailV2,
	}

	c.wsClient.Connect(ctx)

	return c.wsClient.Subscribe(ctx, data)
}
