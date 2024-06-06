//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client.go -destination client_mock.go -package=transactions

package transactions

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
)

const (
	RequestDataType = "timelineTransactions"
)

type ClientInterface interface {
	Get() ([]ResponseItem, error)
}

type Client struct {
	reader portfolio.ReaderInterface
}

func NewClient(reader portfolio.ReaderInterface) Client {
	return Client{
		reader: reader,
	}
}

func (c Client) Get() ([]ResponseItem, error) {
	var result []ResponseItem

	resp, err := c.request("")
	if err != nil {
		return nil, err
	}

	result = slices.Concat(result, resp.Items)

	for resp.Cursors.After != "" {
		resp, err = c.request(resp.Cursors.After)
		if err != nil {
			return nil, err
		}

		result = slices.Concat(result, resp.Items)
	}

	return result, nil
}

func (c Client) request(after string) (websocket.CollectionResponse[ResponseItem], error) {
	var resp websocket.CollectionResponse[ResponseItem]

	msg, err := c.reader.Read(RequestDataType, map[string]any{"after": after})
	if err != nil {
		return resp, fmt.Errorf("could not fetch %s: %w", RequestDataType, err)
	}

	if err = json.Unmarshal(msg.Data(), &resp); err != nil {
		return resp, fmt.Errorf("could not unmarshal %s response: %w", RequestDataType, err)
	}

	return resp, nil
}
