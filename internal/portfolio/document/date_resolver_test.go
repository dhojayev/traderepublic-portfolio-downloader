package document_test

import (
	"fmt"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDateResolver_Resolve(t *testing.T) {
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
	dateResolver := document.NewDateResolver(logger)

	for i, testCase := range testCases {
		for _, doc := range testCase.Transaction.Documents {
			actual, err := dateResolver.Resolve(testCase.Transaction.Timestamp, doc.Detail)

			assert.NoError(t, err, fmt.Sprintf("case %d", i))
			assert.Equal(t, doc.Timestamp, actual)
		}
	}
}
