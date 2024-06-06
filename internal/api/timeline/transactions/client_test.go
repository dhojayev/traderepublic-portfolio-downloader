package transactions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestClient_Get(t *testing.T) {
	t.Parallel()

	testCases := []tests.TestCase{
		fakes.CardSuccessfulTransaction01,
	}

	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	client := transactions.NewClient(readerMock)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineTransactions", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				return filesystem.NewOutputData([]byte(testCase.TimelineTransactionsData.Raw)), nil
			})

		actual, err := client.Get()

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Len(t, actual, 1, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.TimelineTransactionsData.Unmarshalled, actual[0], fmt.Sprintf("case %d", i))
	}
}
