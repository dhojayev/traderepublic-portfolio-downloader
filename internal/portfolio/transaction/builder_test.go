package transaction_test

import (
	"os"
	"testing"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
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

	expectedTime, _ := time.Parse(internal.DefaultTimeFormat, "2022-03-29T09:43:31.570+0000")

	testCases := []struct {
		filepath string
		expected transaction.Model
	}{
		{
			filepath: "../../../tests/data/transaction-details/securities-settlement-variant-1.json",
			expected: transaction.Model{
				UUID: "b20e367c-5542-4fab-9fd6-6faa5e7ab582",
				Instrument: transaction.Instrument{
					ISIN: "DE000SH0MW59",
					Name: "CAC",
					Icon: "logos/FR0003500008/v2",
				},
				Documents: []document.Model{
					{
						ID:    "46e92aa7-df44-4a69-957c-183459753e66",
						URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/",
						Date:  "29.03.2022",
						Title: "Abrechnung",
					},
					{
						ID:    "3c4ccef3-249d-4d10-a54a-18a82fb9475a",
						URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/",
						Date:  "29.03.2022",
						Title: "Kosteninformation",
					},
				},
				Type:       transaction.TypePurchase,
				Timestamp:  expectedTime,
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
	resolver := transaction.NewTypeResolver(log.New())
	builder := transaction.NewBuilder(resolver, log.New())

	for _, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				fileContents, err := os.ReadFile(testCase.filepath)

				return content{data: fileContents}, err
			})

		transactionDetails, err := detailsClient.Get("b20e367c-5542-4fab-9fd6-6faa5e7ab582")
		assert.Nil(t, err)

		actual, err := builder.FromResponse(transactionDetails)
		assert.Nil(t, err)
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestBuilder_BuildDocuments(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filepath string
		expected []document.Model
	}{
		{
			filepath: "../../../tests/data/transaction-details/documents-section-variant-1.json",
			expected: []document.Model{
				{
					ID:    "58acfbab-45fe-4be1-8ec3-3901a6eabf36",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/2023/12/11/138c562b/pb3204932049320940329402394032.pdf",
					Date:  "22.11.2023",
					Title: "Abrechnung",
				},
				{
					ID:    "3076d454-edcc-4653-a170-31bcd06164c1",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/2023/12/11/12mn31292/pb32492394319490103012312.pdf",
					Date:  "23.11.2023",
					Title: "Kosteninformation",
				},
			},
		},
		{
			filepath: "../../../tests/data/transaction-details/documents-section-variant-2.json",
			expected: []document.Model{
				{
					ID:    "51f4e1cf-30ac-4c6b-92cb-afb5bba19e20",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/2024/3/25/9234rjd23/pb12390210938921839218123012.pdf",
					Title: "Abrechnung Ausführung",
				},
				{
					ID:    "3fc1e7e6-2fa9-43bf-af6e-6f8e9f744226",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/2024/3/25/234391d/pb1234991203901249023902.pdf",
					Title: "Kosteninformation",
				},
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

		actual, err := builder.BuildDocuments(transactionDetails)
		assert.Nil(t, err)
		assert.Len(t, actual, len(testCase.expected))
		assert.Equal(t, testCase.expected, actual)
	}
}
