package transactions

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
)

const (
	dataType = "timelineTransactions"
)

type Client struct {
	reader portfolio.ReaderInterface
}

func NewClient(reader portfolio.ReaderInterface) Client {
	return Client{
		reader: reader,
	}
}

func (c *Client) Get() ([]TransactionResponse, error) {
	var result []TransactionResponse

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

func (c *Client) request(after string) (websocket.CollectionResponse[TransactionResponse], error) {
	var resp websocket.CollectionResponse[TransactionResponse]

	msg, err := c.reader.Read(dataType, map[string]any{"after": after})
	if err != nil {
		return resp, fmt.Errorf("could not fetch %s: %w", dataType, err)
	}

	if err = json.Unmarshal(msg.Data(), &resp); err != nil {
		return resp, fmt.Errorf("could not unmarshal %s response: %w", dataType, err)
	}

	return resp, nil
}
