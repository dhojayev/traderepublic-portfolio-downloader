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
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestModelBuilderBuildSupported(t *testing.T) {
	t.Parallel()

	testCases := []tests.TestCase{
		fakes.BenefitsSpareChangeExecution01,
		fakes.BenefitsSavebackExecution01,
		fakes.Credit01,
		fakes.OrderExecuted01,
		fakes.OrderExecuted02,
		fakes.OrderExecuted03,
		fakes.PaymentInbound01,
		fakes.PaymentInboundSepaDirectDebit01,
		fakes.PaymentOutbound01,
		fakes.SavingsPlanExecuted01,
	}

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	detailsClient := details.NewClient(readerMock)
	resolver := details.NewTypeResolver(logger)
	documentDateResolver := document.NewDateResolver(logger)
	documentBuilder := document.NewModelBuilder(documentDateResolver, logger)
	builderFactory := transaction.NewModelBuilderFactory(resolver, documentBuilder, logger)

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

		if err != nil {
			return
		}

		actual, err := builder.Build()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.Transaction, actual, fmt.Sprintf("case %d", i))
	}
}

func TestModelBuilderBuildUnsupported(t *testing.T) {
	t.Parallel()

	testCases := []tests.TestCase{
		fakes.InterestPayoutCreated01,
	}

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	detailsClient := details.NewClient(readerMock)
	resolver := details.NewTypeResolver(logger)
	documentDateResolver := document.NewDateResolver(logger)
	documentBuilder := document.NewModelBuilder(documentDateResolver, logger)
	builderFactory := transaction.NewModelBuilderFactory(resolver, documentBuilder, logger)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				return filesystem.NewOutputData([]byte(testCase.ResponseJSON)), nil
			})

		response, err := detailsClient.Get("b20e367c-5542-4fab-9fd6-6faa5e7ab582")
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		_, err = builderFactory.Create(testCase.EventType, response)
		assert.Error(t, err, fmt.Sprintf("case %d", i))
	}
}
