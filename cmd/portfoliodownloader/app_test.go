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
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestItDoesReturnErrorIfTransactionDetailsCannotBeFetched(t *testing.T) {
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
		Read(details.RequestDataType, map[string]any{"id": "0e5cf3cb-0f4d-4905-ae5f-ec0a530de6ca"}).
		Return(filesystem.OutputData{}, fmt.Errorf("could not fetch %s: %w", details.RequestDataType, websocket.ErrMsgErrorStateReceived))

	err := app.Run()

	assert.NoError(t, err)
}

func TestItDoesReturnErrorIfTransactionTypeUnsupported(t *testing.T) {
	t.Parallel()

	testCases := make([]fakes.TestCase, 0)
	testCases = append(testCases, fakes.TestCasesUnsupported...)
	testCases = append(testCases, fakes.TestCasesUnknown...)

	ctrl := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(ctrl)
	transactionRepoMock := transaction.NewMockRepositoryInterface(ctrl)
	csvReaderMock := filesystem.NewMockCSVReaderInterface(ctrl)
	csvWriterMock := filesystem.NewMockCSVWriterInterface(ctrl)
	documentDownloaderMock := document.NewMockDownloaderInterface(ctrl)

	logger := log.New()
	transactionsClient := transactions.NewClient(readerMock)
	transactionsTypeResolver := transactions.NewEventTypeResolver(logger)
	detailsClient := details.NewClient(readerMock)
	detailsTypeResolver := details.NewTypeResolver(logger)
	documentDateResolver := document.NewDateResolver(logger)
	documentBuilder := document.NewModelBuilder(documentDateResolver, logger)
	builderFactory := transaction.NewModelBuilderFactory(detailsTypeResolver, documentBuilder, logger)
	csvEntryFactory := transaction.NewCSVEntryFactory(logger)
	processor := transaction.NewProcessor(builderFactory, transactionRepoMock, csvEntryFactory, csvReaderMock, csvWriterMock, documentDownloaderMock, logger)
	app := portfoliodownloader.NewApp(transactionsClient, transactionsTypeResolver, detailsClient, processor, logger)

	csvReaderMock.EXPECT().Read(gomock.Any()).AnyTimes().Return([]filesystem.CSVEntry{}, nil)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read(transactions.RequestDataType, gomock.Any()).
			Times(1).
			Return(filesystem.NewOutputData([]byte(testCase.TimelineTransactionsData.Raw)), nil)

		readerMock.
			EXPECT().
			Read(details.RequestDataType, map[string]any{"id": testCase.TimelineTransactionsData.Unmarshalled.Action.Payload}).
			Times(1).
			Return(filesystem.NewOutputData([]byte(testCase.TimelineDetailsData.Raw)), nil)

		err := app.Run()

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
	}
}
