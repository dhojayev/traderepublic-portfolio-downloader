//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client.go -destination client_mock.go -package=details

package details

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	log "github.com/sirupsen/logrus"
)

const (
	RequestDataType = "timelineDetailV2"
)

type ClientInterface interface {
	api.WSDetailsGetterClientInterface
}

type Client struct {
	api.WSClient
}

func NewClient(reader reader.Interface, logger *log.Logger) Client {
	return Client{api.NewWSClient(RequestDataType, reader, logger)}
}
