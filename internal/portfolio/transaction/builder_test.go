package transaction_test

import (
	"os"
	"testing"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type content struct {
	data []byte
}

func (c content) Data() []byte {
	return c.data
}

func TestBuildPurchaseTransactions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filepath string
		expected transaction.Model
	}{
		{
			filepath: "../../../tests/data/transaction-details/securities-settlement-variant-1.json",
			expected: transaction.Model{
				UUID:           "b20e367c-5542-4fab-9fd6-6faa5e7ab582",
				InstrumentISIN: "DE000SH0MW59",
				Instrument: transaction.Instrument{
					ISIN: "DE000SH0MW59",
					Name: "CAC",
					Icon: "logos/FR0003500008/v2",
				},
				Documents:  []transaction.Document{},
				Type:       transaction.TypePurchase,
				Timestamp:  time.Now(),
				Status:     "executed",
				Yield:      0,
				Profit:     0,
				Shares:     40,
				Rate:       9.87,
				Commission: 1,
				Total:      395.80,
			},
		},
	}

	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	detailsClient := details.NewClient(readerMock)
	resolverMock := transaction.NewMockTypeResolverInterface(controller)
	builder := transaction.NewBuilder(resolverMock, log.New())

	for _, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				fileContents, err := os.ReadFile(testCase.filepath)

				return content{data: fileContents}, err
			})

		transactionDetails, err := detailsClient.Get("2d7c03e4-15f9-4427-88d2-0586c5b057d2")
		assert.Nil(t, err)

		documents, err := builder.BuildDocuments(transactionDetails)
		assert.Nil(t, err)
		assert.Len(t, documents, len(testCase.expected))
		assert.Equal(t, testCase.expected, documents)
	}
}
