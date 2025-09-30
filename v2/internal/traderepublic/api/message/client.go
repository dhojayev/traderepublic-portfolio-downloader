package message

import (
	"context"
	"log/slog"
	"strconv"
	"sync"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type ClientInterface interface {
	SubscribeToTimelineTransactions(ctx context.Context) error
	SubscribeToTimelineDetailV2(ctx context.Context, itemID string) error
}

type Client struct {
	eventBus           *bus.EventBus
	credentialsService auth.CredentialsServiceInterface
	wsClient           websocketclient.ClientInterface
}

func NewClient(eventBus *bus.EventBus, credentialsService auth.CredentialsServiceInterface, wsClient websocketclient.ClientInterface) *Client {
	return &Client{
		eventBus:           eventBus,
		credentialsService: credentialsService,
		wsClient:           wsClient,
	}
}

// SubscribeToTimelineTransactions subscribes to timeline transactions data.
func (c *Client) SubscribeToTimelineTransactions(ctx context.Context) error {
	ch, err := c.SubscribeToTimelineTransactionsWithCursor(ctx, "")
	if err != nil {
		return err
	}

	data := <-ch
	counter := int64(1)

	c.eventBus.Publish(bus.NewEvent(
		bus.TopicTimelineTransactions,
		strconv.FormatInt(counter, 10),
		bus.EventNameTimelineTransactionsReceived,
		data,
	))

	var response traderepublic.TimelineTransactionsSchemaJson

	err = response.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	var mu sync.Mutex

	go func() {
		for response.Cursors.After != nil {
			ch, err = c.SubscribeToTimelineTransactionsWithCursor(ctx, *response.Cursors.After)
			if err != nil {
				slog.Error("error subscribing to timeline transactions", "error", err)

				return
			}

			data = <-ch

			mu.Lock()

			counter++

			mu.Unlock()

			c.eventBus.Publish(bus.NewEvent(
				bus.TopicTimelineTransactions,
				strconv.FormatInt(counter, 10),
				bus.EventNameTimelineTransactionsReceived,
				data,
			))

			err = response.UnmarshalJSON(data)
			if err != nil {
				slog.Error("error subscribing to timeline transactions", "error", err)

				return
			}
		}
	}()

	return nil
}

// SubscribeToTimelineTransactionsWithCursor subscribes to timeline transactions data with a cursor.
func (c *Client) SubscribeToTimelineTransactionsWithCursor(ctx context.Context, cursor string) (<-chan []byte, error) {
	data := traderepublic.WebsocketSubRequestSchemaJson{
		Token: c.credentialsService.GetToken().Session(),
		Type:  traderepublic.WebsocketSubRequestSchemaJsonTypeTimelineTransactions,
		After: &cursor,
	}

	return c.wsClient.Subscribe(data)
}

// SubscribeToTimelineDetail subscribes to timeline detail data.
func (c *Client) SubscribeToTimelineDetailV2(ctx context.Context, itemID string) error {
	data := traderepublic.WebsocketSubRequestSchemaJson{
		Id:    &itemID,
		Token: c.credentialsService.GetToken().Session(),
		Type:  traderepublic.WebsocketSubRequestSchemaJsonTypeTimelineDetailV2,
	}

	ch, err := c.wsClient.Subscribe(data)
	if err != nil {
		return nil
	}

	go func() {
		data := <-ch

		c.eventBus.Publish(bus.NewEvent(
			bus.TopicTimelineDetailsV2,
			itemID,
			bus.EventNameTimelineDetailV2Received,
			data,
		))
	}()

	return nil
}
