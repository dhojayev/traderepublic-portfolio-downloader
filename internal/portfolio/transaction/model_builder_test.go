package transaction_test

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestPurchaseBuilderBuild(t *testing.T) {
	t.Parallel()

	testCases := []tests.TestCase{
		fakes.OrderExecuted01,
		fakes.BenefitsSpareChangeExecution01,
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
				return filesystem.NewOutputData([]byte(testCase.ResponseJSON)), nil
			})

		response, err := detailsClient.Get("b20e367c-5542-4fab-9fd6-6faa5e7ab582")
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		builder, err := builderFactory.Create(testCase.EventType, response)
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		actual, err := builder.Build()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.Transaction, actual, fmt.Sprintf("case %d", i))
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
