//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client.go -destination client_mock.go -package=activitylog

package activitylog

import (
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api"
)

const RequestDataType = "timelineActivityLog"

type ClientInterface interface {
	api.WSListGetterClientInterface
}

type Client struct {
	api.WSClient
}

func NewClient(reader reader.Interface, logger *log.Logger) Client {
	return Client{api.NewWSClient(RequestDataType, reader, logger)}
}
