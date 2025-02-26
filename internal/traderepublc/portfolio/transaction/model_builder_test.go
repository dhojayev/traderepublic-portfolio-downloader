package transaction_test

import (
	"fmt"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
	details_test "github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes/details"
)

func TestModelBuilderBuildSupported(t *testing.T) {
	t.Parallel()

	testCases := details_test.TestCasesSupported

	if len(testCases) == 0 {
		t.Error("no test cases found")
	}

	logger := log.New()
	logger.Out = io.Discard

	controller := gomock.NewController(t)
	readerMock := reader.NewMockInterface(controller)
	detailsClient := details.NewClient(readerMock, logger)
	normalizer := details.NewTransactionResponseNormalizer(logger)
	builderFactory := transaction.ProvideModelBuilderFactory(logger)

	for testCaseName, testCase := range testCases {
		if testCase.Parent == nil {
			t.Errorf("case '%s' does not contain parent", testCaseName)

			continue
		}

		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (reader.JSONResponse, error) {
				return reader.NewJSONResponse(testCase.RawResponse), nil
			})

		var response details.Response

		err := detailsClient.Details("b20e367c-5542-4fab-9fd6-6faa5e7ab582", &response)
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		normalizedResponse, err := normalizer.Normalize(response)
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		builder, err := builderFactory.Create(transactions.EventType(testCase.Parent.Unmarshalled.EventType), normalizedResponse)
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		if err != nil {
			return
		}

		actual, err := builder.Build()
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))
		assert.Equal(t, testCase.Model, actual, fmt.Sprintf("case '%s'", testCaseName))
	}
}

func TestModelBuilderBuildUnsupported(t *testing.T) {
	t.Parallel()

	testCases := details_test.TestCasesUnsupported

	if len(testCases) == 0 {
		t.Error("no test cases found")
	}

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := reader.NewMockInterface(controller)
	detailsClient := details.NewClient(readerMock, logger)
	normalizer := details.NewTransactionResponseNormalizer(logger)
	builderFactory := transaction.ProvideModelBuilderFactory(logger)

	for testCaseName, testCase := range testCases {
		if testCase.Parent == nil {
			t.Errorf("case '%s' does not contain parent", testCaseName)

			continue
		}

		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (reader.JSONResponse, error) {
				return reader.NewJSONResponse(testCase.RawResponse), nil
			})

		var response details.Response

		err := detailsClient.Details("b20e367c-5542-4fab-9fd6-6faa5e7ab582", &response)
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		normalizedResponse, _ := normalizer.Normalize(response)

		_, err = builderFactory.Create(transactions.EventType(testCase.Parent.Unmarshalled.EventType), normalizedResponse)
		assert.Error(t, err, fmt.Sprintf("case '%s'", testCaseName))
	}
}

func TestModelBuilderBuildUnknown(t *testing.T) {
	t.Parallel()

	testCases := details_test.TestCasesUnknown

	if len(testCases) == 0 {
		t.Error("no test cases found")
	}

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := reader.NewMockInterface(controller)
	detailsClient := details.NewClient(readerMock, logger)
	normalizer := details.NewTransactionResponseNormalizer(logger)
	builderFactory := transaction.ProvideModelBuilderFactory(logger)

	for testCaseName, testCase := range testCases {
		if testCase.Parent == nil {
			t.Errorf("case '%s' does not contain parent", testCaseName)

			continue
		}

		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (reader.JSONResponse, error) {
				return reader.NewJSONResponse(testCase.RawResponse), nil
			})

		var response details.Response

		err := detailsClient.Details("b20e367c-5542-4fab-9fd6-6faa5e7ab582", &response)
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		normalizedResponse, err := normalizer.Normalize(response)
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		builder, err := builderFactory.Create(transactions.EventType(testCase.Parent.Unmarshalled.EventType), normalizedResponse)
		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		if err != nil {
			return
		}

		_, err = builder.Build()
		assert.ErrorIs(t, err, transaction.ErrModelBuilderInsufficientDataResolved, fmt.Sprintf("case '%s'", testCaseName))
	}
}
