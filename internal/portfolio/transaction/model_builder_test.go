package transaction_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

func TestPurchaseBuilderBuild(t *testing.T) {
	t.Parallel()

	expectedTimes := make([]time.Time, 2)
	expectedTimes[0], _ = time.Parse(internal.DefaultTimeFormat, "2022-03-29T09:43:31.570+0000")
	expectedTimes[1], _ = time.Parse(internal.DefaultTimeFormat, "2024-01-04T12:26:52.110+0000")

	testCases := []struct {
		filepath  string
		eventType transactions.EventType
		expected  transaction.Model
	}{
		{
			filepath:  "../../../tests/data/transaction-details/order-executed-04.json",
			eventType: transactions.EventTypeOrderExecuted,
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
				Timestamp:  expectedTimes[0],
				Status:     "executed",
				Shares:     40,
				Rate:       9.87,
				Commission: 1,
				Total:      395.80,
			},
		},
		{
			filepath:  "../../../tests/data/transaction-details/benefits-spare-change-execution-01.json",
			eventType: transactions.EventTypeBenefitsSpareChangeExecution,
			expected: transaction.Model{
				UUID: "265cb9c0-664a-45d4-b179-3061f196dd2a",
				Instrument: transaction.Instrument{
					ISIN: "DE000A0F5UF5",
					Name: "NASDAQ100 USD (Dist)",
					Icon: "logos/DE000A0F5UF5/v2",
				},
				Documents: []document.Model{
					{
						ID:    "9df4c2e1-0de2-4900-aa8c-af5371ed58f6",
						URL:   "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
						Title: "Deaktivierung",
					},
					{
						ID:    "3a8ebf86-a2bb-463e-8bfd-28fd705359ff",
						URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
						Title: "Abrechnung Ausführung",
					},
					{
						ID:    "e2dfa755-e039-45c7-b7bb-1ac024844f75",
						URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
						Title: "Kosteninformation",
					},
				},
				Type:      transaction.TypeRoundUp,
				Timestamp: expectedTimes[1],
				Status:    "executed",
				Shares:    0.006882,
				Rate:      158.38,
				Total:     1.09,
			},
		},
	}

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	detailsClient := details.NewClient(readerMock)
	resolver := details.NewTypeResolver(logger)
	builderFactory := transaction.NewModelBuilderFactory(resolver, logger)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				fileContents, err := os.ReadFile(testCase.filepath)

				return filesystem.NewOutputData(fileContents), err
			})

		response, err := detailsClient.Get("b20e367c-5542-4fab-9fd6-6faa5e7ab582")
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		builder, err := builderFactory.Create(testCase.eventType, response)
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		actual, err := builder.Build()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.expected, actual, fmt.Sprintf("case %d", i))
	}
}

func TestPurchaseBuilderBuildDocuments(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filepath  string
		eventType transactions.EventType
		expected  []document.Model
	}{
		{
			filepath:  "../../../tests/data/transaction-details/order-executed-03.json",
			eventType: transactions.EventTypeOrderExecuted,
			expected: []document.Model{
				{
					ID:    "f17b2237-0e32-410e-b38b-8638600ffbb0",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
					Date:  "11.03.2024",
					Title: "Abrechnung",
				},
				{
					ID:    "3c214355-dc5a-488a-b780-b28fb66b66c8",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
					Date:  "27.02.2024",
					Title: "Auftragsbestätigung",
				},
				{
					ID:    "21a13acc-7f3c-4156-8365-be8089006ac4",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
					Date:  "12.02.2024",
					Title: "Kosteninformation",
				},
			},
		},
		{
			filepath:  "../../../tests/data/transaction-details/benefits-spare-change-execution-01.json",
			eventType: transactions.EventTypeBenefitsSpareChangeExecution,
			expected: []document.Model{
				{
					ID:    "9df4c2e1-0de2-4900-aa8c-af5371ed58f6",
					URL:   "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
					Title: "Deaktivierung",
				},
				{
					ID:    "3a8ebf86-a2bb-463e-8bfd-28fd705359ff",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
					Title: "Abrechnung Ausführung",
				},
				{
					ID:    "e2dfa755-e039-45c7-b7bb-1ac024844f75",
					URL:   "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
					Title: "Kosteninformation",
				},
			},
		},
	}

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	detailsClient := details.NewClient(readerMock)
	resolver := details.NewTypeResolver(logger)
	builderFactory := transaction.NewModelBuilderFactory(resolver, logger)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				fileContents, err := os.ReadFile(testCase.filepath)

				return filesystem.NewOutputData(fileContents), err
			})

		response, err := detailsClient.Get("2d7c03e4-15f9-4427-88d2-0586c5b057d2")
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		builder, err := builderFactory.Create(testCase.eventType, response)
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		model, err := builder.Build()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		actual := model.Documents
		assert.Len(t, actual, len(testCase.expected), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.expected, actual, fmt.Sprintf("case %d", i))
	}
}
