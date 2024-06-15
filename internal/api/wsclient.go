package api

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	log "github.com/sirupsen/logrus"
)

// WSListGetterClientInterface represents a Websocket client that can make "list" requests.
type WSListGetterClientInterface interface {
	List(v any) error
}

// WSDetailsGetterClientInterface represents a Websocket client that can make "details" requests.
type WSDetailsGetterClientInterface interface {
	Details(itemID string, v any) error
}

// WSClient Websocket client.
type WSClient struct {
	dataType string
	reader   portfolio.ReaderInterface
	logger   *log.Logger
}

func NewWSClient(dataType string, reader portfolio.ReaderInterface, logger *log.Logger) WSClient {
	return WSClient{
		dataType: dataType,
		reader:   reader,
		logger:   logger,
	}
}

func (c WSClient) List(v any) error {
	var items []any

	var resp WSListResponse

	err := c.request(&resp, nil)
	if err != nil {
		return err
	}

	items = slices.Concat(items, resp.Items)

	for resp.Cursors.After != "" {
		cursorAfter := resp.Cursors.After
		resp = WSListResponse{}

		err = c.request(&resp, map[string]any{"after": cursorAfter})
		if err != nil {
			return err
		}

		items = slices.Concat(items, resp.Items)
	}

	itemsMarshalled, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("could not marshal %s list items: %w", c.dataType, err)
	}

	if err = json.Unmarshal(itemsMarshalled, v); err != nil {
		return fmt.Errorf("could not unmarshal %s list response: %w", c.dataType, err)
	}

	return nil
}

func (c WSClient) Details(itemID string, v any) error {
	if err := c.request(v, map[string]any{"id": itemID}); err != nil {
		return err
	}

	return nil
}

func (c WSClient) request(v any, data map[string]any) error {
	msg, err := c.reader.Read(c.dataType, data)
	if err != nil {
		return fmt.Errorf("could not fetch %s: %w", c.dataType, err)
	}

	if err = json.Unmarshal(msg.Data(), v); err != nil {
		return fmt.Errorf("could not unmarshal %s response: %w", c.dataType, err)
	}

	return nil
}
