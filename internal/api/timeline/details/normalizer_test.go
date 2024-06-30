package details_test

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestItNormalizesSupportTransactionDetails(t *testing.T) {
	t.Parallel()

	testCases := map[string]fakes.TransactionTestCase{
		"PaymentInbound01":                fakes.PaymentInbound01,
		"PaymentInboundSepaDirectDebit01": fakes.PaymentInboundSepaDirectDebit01,
		"InterestPayoutCreated01":         fakes.InterestPayoutCreated01,
		"SavingsPlanExecuted01":           fakes.SavingsPlanExecuted01,
		"OrderExecuted02":                 fakes.OrderExecuted02,
		"Credit01":                        fakes.Credit01,
		"BenefitsSpareChangeExecution01":  fakes.BenefitsSpareChangeExecution01,
	}

	logger := log.New()
	logger.Out = io.Discard

	normalizer := details.NewTransactionResponseNormalizer(logger)

	for testCaseName, testCase := range testCases {
		var response details.Response

		err := json.Unmarshal(testCase.TimelineDetailsData.Raw, &response)

		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		actual, err := normalizer.Normalize(response)

		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))
		assert.Equal(t, testCase.TimelineDetailsData.Normalized, actual, fmt.Sprintf("case '%s'", testCaseName))
	}
}
