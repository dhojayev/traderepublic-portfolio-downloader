package transactions_test

import (
	"fmt"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestClient_Get(t *testing.T) {
	t.Parallel()

	testCases := map[string]fakes.TransactionTestCase{
		"CardSuccessfulTransaction01": fakes.CardSuccessfulTransaction01,
		"CardSuccessfulTransaction02": fakes.CardSuccessfulTransaction02,
	}

	logger := log.New()
	logger.Out = io.Discard

	controller := gomock.NewController(t)
	readerMock := reader.NewMockInterface(controller)
	client := transactions.NewClient(readerMock, logger)

	for testCaseName, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineTransactions", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (reader.JSONResponse, error) {
				return reader.NewJSONResponse(testCase.TimelineTransactionsData.Raw), nil
			})

		var actual []transactions.ResponseItem
		err := client.List(&actual)

		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		if err != nil {
			continue
		}

		assert.Len(t, actual, 1, fmt.Sprintf("case '%s'", testCaseName))
		assert.Equal(t, testCase.TimelineTransactionsData.Unmarshalled, actual[0], fmt.Sprintf("case '%s'", testCaseName))
	}
}
