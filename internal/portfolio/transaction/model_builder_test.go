package transaction_test

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
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

	testCases := []tests.TestCase{
		fakes.OrderExecuted03,
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

		response, err := detailsClient.Get("2d7c03e4-15f9-4427-88d2-0586c5b057d2")
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		builder, err := builderFactory.Create(testCase.EventType, response)
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		model, err := builder.Build()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		actual := model.Documents
		assert.Len(t, actual, len(testCase.Transaction.Documents), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.Transaction.Documents, actual, fmt.Sprintf("case %d", i))
	}
}
