package portfoliodownloader_test

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestItDoesReturnErrorIfTransactionDetailsCannotBeFetched(t *testing.T) {
	t.Parallel()

	logger := log.New()
	ctrl := gomock.NewController(t)
	activityLogClientMock := activitylog.NewMockClientInterface(ctrl)
	transactionsClientMock := transactions.NewMockClientInterface(ctrl)
	detailsReaderMock := reader.NewMockInterface(ctrl)
	typeResolverMock := transactions.NewMockEventTypeResolverInterface(ctrl)
	processorMock := transaction.NewMockProcessorInterface(ctrl)

	detailsClient := details.NewClient(detailsReaderMock, logger)
	app := portfoliodownloader.NewApp(transactionsClientMock, typeResolverMock, detailsClient, processorMock, activityLogClientMock, logger)

	var responses []transactions.ResponseItem

	activityLogClientMock.EXPECT().List(gomock.Any()).Times(1).Return(nil)

	transactionsClientMock.
		EXPECT().
		List(&responses).
		SetArg(0, []transactions.ResponseItem{
			{
				Action: transactions.ResponseItemAction{
					Payload: "0e5cf3cb-0f4d-4905-ae5f-ec0a530de6ca",
					Type:    "timelineDetail",
				},
				ID: "0e5cf3cb-0f4d-4905-ae5f-ec0a530de6ca",
			},
		})

	detailsReaderMock.
		EXPECT().
		Read(details.RequestDataType, map[string]any{"id": "0e5cf3cb-0f4d-4905-ae5f-ec0a530de6ca"}).
		Return(reader.NewJSONResponse(nil), fmt.Errorf("could not fetch %s: %w", details.RequestDataType, websocket.ErrMsgErrorStateReceived))

	err := app.Run()

	assert.NoError(t, err)
}

func TestItDoesReturnErrorIfTransactionTypeUnsupported(t *testing.T) {
	t.Parallel()

	testCases := make([]fakes.TransactionTestCase, 0)
	testCases = append(testCases, fakes.TransactionTestCasesUnsupported...)
	testCases = append(testCases, fakes.TransactionTestCasesUnknown...)

	ctrl := gomock.NewController(t)
	readerMock := reader.NewMockInterface(ctrl)
	transactionRepoMock := transaction.NewMockRepositoryInterface(ctrl)
	csvReaderMock := filesystem.NewMockCSVReaderInterface(ctrl)
	csvWriterMock := filesystem.NewMockCSVWriterInterface(ctrl)
	documentDownloaderMock := document.NewMockDownloaderInterface(ctrl)

	logger := log.New()
	activityLogClientMock := activitylog.NewMockClientInterface(ctrl)
	transactionsClient := transactions.NewClient(readerMock, logger)
	transactionsTypeResolver := transactions.NewEventTypeResolver(logger)
	detailsClient := details.NewClient(readerMock, logger)
	detailsTypeResolver := details.NewTypeResolver(logger)
	documentDateResolver := document.NewDateResolver(logger)
	documentBuilder := document.NewModelBuilder(documentDateResolver, logger)
	builderFactory := transaction.NewModelBuilderFactory(detailsTypeResolver, documentBuilder, logger)
	csvEntryFactory := transaction.NewCSVEntryFactory(logger)
	processor := transaction.NewProcessor(builderFactory, transactionRepoMock, csvEntryFactory, csvReaderMock, csvWriterMock, documentDownloaderMock, logger)
	app := portfoliodownloader.NewApp(transactionsClient, transactionsTypeResolver, detailsClient, processor, activityLogClientMock, logger)

	csvReaderMock.EXPECT().Read(gomock.Any()).AnyTimes().Return([]filesystem.CSVEntry{}, nil)

	for i, testCase := range testCases {
		activityLogClientMock.EXPECT().List(gomock.Any()).Times(1).Return(nil)

		readerMock.
			EXPECT().
			Read(transactions.RequestDataType, gomock.Any()).
			Times(1).
			Return(reader.NewJSONResponse(testCase.TimelineTransactionsData.Raw), nil)

		readerMock.
			EXPECT().
			Read(details.RequestDataType, map[string]any{"id": testCase.TimelineTransactionsData.Unmarshalled.Action.Payload}).
			Times(1).
			Return(reader.NewJSONResponse(testCase.TimelineDetailsData.Raw), nil)

		err := app.Run()

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
	}
}
