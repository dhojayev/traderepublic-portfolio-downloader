package details_test

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	details_test "github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes/details"
)

func TestItNormalizesSupportTransactionDetails(t *testing.T) {
	t.Parallel()

	testCases := details_test.TestCasesSupported

	if len(testCases) == 0 {
		t.Error("no test cases found")
	}

	logger := log.New()
	logger.Out = io.Discard

	normalizer := details.NewTransactionResponseNormalizer(logger)

	for testCaseName, testCase := range testCases {
		if reflect.DeepEqual(testCase.Unmarshalled, details.NormalizedResponse{}) {
			t.Logf("case '%s' does not contain unmarshalled response", testCaseName)

			continue
		}

		var response details.Response

		err := json.Unmarshal(testCase.RawResponse, &response)

		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		actual, err := normalizer.Normalize(response)

		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))
		assert.Equal(t, testCase.Unmarshalled, actual, fmt.Sprintf("case '%s'", testCaseName))
	}
}
