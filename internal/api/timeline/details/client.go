//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client.go -destination client_mock.go -package=details

package details

import (
	"encoding/json"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
)

const (
	RequestDataType = "timelineDetailV2"
)

type ClientInterface interface {
	Get(eventID string) (Response, error)
}

type Client struct {
	reader portfolio.ReaderInterface
}

func NewClient(reader portfolio.ReaderInterface) Client {
	return Client{
		reader: reader,
	}
}

func (c Client) Get(eventID string) (Response, error) {
	var response Response

	msg, err := c.reader.Read(RequestDataType, map[string]any{"id": eventID})
	if err != nil {
		return response, fmt.Errorf("could not fetch %s: %w", RequestDataType, err)
	}

	if err = json.Unmarshal(msg.Data(), &response); err != nil {
		return response, fmt.Errorf("could not unmarshal %s response: %w", RequestDataType, err)
	}

	return response, nil
}
