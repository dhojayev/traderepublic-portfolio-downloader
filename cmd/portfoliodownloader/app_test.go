package portfoliodownloader_test

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

func TestApp_Run(t *testing.T) {
	t.Parallel()

	logger := log.New()
	ctrl := gomock.NewController(t)
	transactionsClientMock := transactions.NewMockClientInterface(ctrl)
	detailsReaderMock := portfolio.NewMockReaderInterface(ctrl)
	typeResolverMock := transactions.NewMockEventTypeResolverInterface(ctrl)
	processorMock := transaction.NewMockProcessorInterface(ctrl)

	detailsClient := details.NewClient(detailsReaderMock)
	app := portfoliodownloader.NewApp(transactionsClientMock, typeResolverMock, detailsClient, processorMock, logger)

	transactionsClientMock.
		EXPECT().
		Get().
		Return([]transactions.ResponseItem{
			{
				Action: transactions.ResponseItemAction{
					Payload: "0e5cf3cb-0f4d-4905-ae5f-ec0a530de6ca",
					Type:    "timelineDetail",
				},
				ID: "0e5cf3cb-0f4d-4905-ae5f-ec0a530de6ca",
			},
		}, nil)

	detailsReaderMock.
		EXPECT().
		Read(details.DataType, map[string]any{"id": "0e5cf3cb-0f4d-4905-ae5f-ec0a530de6ca"}).
		Return(filesystem.OutputData{}, fmt.Errorf("could not fetch %s: %w", details.DataType, websocket.ErrErrorStateReceived))

	err := app.Run()

	assert.NoError(t, err)
}
