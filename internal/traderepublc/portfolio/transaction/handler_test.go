package transaction_test

import (
	"fmt"
	"io"
	"maps"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestItDoesReturnErrorIfTransactionDetailsCannotBeFetched(t *testing.T) {
	t.Parallel()

	logger := log.New()
	logger.Out = io.Discard

	ctrl := gomock.NewController(t)
	listClientMock := transactions.NewMockClientInterface(ctrl)
	detailsReaderMock := reader.NewMockInterface(ctrl)
	typeResolverMock := transactions.NewMockEventTypeResolverInterface(ctrl)
	processorMock := transaction.NewMockProcessorInterface(ctrl)
	detailsClient := details.NewClient(detailsReaderMock, logger)
	normalizer := details.NewTransactionResponseNormalizer(logger)
	handler := transaction.NewHandler(listClientMock, detailsClient, normalizer, typeResolverMock, processorMock, logger)

	var responses []transactions.ResponseItem

	listClientMock.
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

	err := handler.Handle()

	assert.NoError(t, err)
}

func TestItDoesReturnErrorIfTransactionTypeUnsupported(t *testing.T) {
	t.Parallel()

	testCases := make(map[string]fakes.TransactionTestCase)
	maps.Copy(testCases, fakes.TransactionTestCasesUnsupported)
	maps.Copy(testCases, fakes.TransactionTestCasesUnknown)

	ctrl := gomock.NewController(t)
	readerMock := reader.NewMockInterface(ctrl)
	transactionRepoMock := transaction.NewMockRepositoryInterface(ctrl)
	csvReaderMock := filesystem.NewMockCSVReaderInterface(ctrl)
	csvWriterMock := filesystem.NewMockCSVWriterInterface(ctrl)
	documentDownloaderMock := document.NewMockDownloaderInterface(ctrl)

	logger := log.New()
	logger.Out = io.Discard

	listClient := transactions.NewClient(readerMock, logger)
	transactionsTypeResolver := transactions.NewEventTypeResolver(logger)
	detailsClient := details.NewClient(readerMock, logger)
	builderFactory := transaction.ProvideModelBuilderFactory(logger)
	csvEntryFactory := transaction.NewCSVEntryFactory(logger)
	processor := transaction.NewProcessor(builderFactory, transactionRepoMock, csvEntryFactory, csvReaderMock, csvWriterMock, documentDownloaderMock, logger)
	normalizer := details.NewTransactionResponseNormalizer(logger)
	handler := transaction.NewHandler(listClient, detailsClient, normalizer, transactionsTypeResolver, processor, logger)

	csvReaderMock.EXPECT().Read(gomock.Any()).AnyTimes().Return([]filesystem.DepotTransactionCSVEntry{}, nil)

	for testCaseName, testCase := range testCases {
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

		err := handler.Handle()

		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))
	}
}
