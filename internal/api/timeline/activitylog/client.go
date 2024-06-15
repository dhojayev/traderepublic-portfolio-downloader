//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client.go -destination client_mock.go -package=activitylog

package activitylog

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	log "github.com/sirupsen/logrus"
)

const RequestDataType = "timelineActivityLog"

type ClientInterface interface {
	api.WSListGetterClientInterface
}

type Client struct {
	api.WSClient
}

func NewClient(reader portfolio.ReaderInterface, logger *log.Logger) Client {
	return Client{api.NewWSClient(RequestDataType, reader, logger)}
}
