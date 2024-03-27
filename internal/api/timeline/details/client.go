package details

import (
	"encoding/json"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
)

const (
	dataType = "timelineDetailV2"
)

type Client struct {
	retriever portfolio.ReaderInterface
}

func NewClient(retriever portfolio.ReaderInterface) Client {
	return Client{
		retriever: retriever,
	}
}

func (c *Client) Get(eventID string) (Response, error) {
	var response Response

	msg, err := c.retriever.Read(dataType, map[string]any{"id": eventID})
	if err != nil {
		return response, fmt.Errorf("could not fetch %s: %w", dataType, err)
	}

	if err = json.Unmarshal(msg.Data(), &response); err != nil {
		return response, fmt.Errorf("could not unmarshal %s response: %w", dataType, err)
	}

	return response, nil
}
