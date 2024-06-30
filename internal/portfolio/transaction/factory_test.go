package transaction_test

import (
	"fmt"
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestMakeSupported(t *testing.T) {
	t.Parallel()

	testCases := fakes.TransactionTestCasesSupported
	factory := transaction.NewCSVEntryFactory(log.New())

	for testCaseName, testCase := range testCases {
		actual, err := factory.Make(testCase.Transaction)

		assert.NoError(t, err, fmt.Sprintf("case '%s'", testCaseName))

		if err != nil {
			continue
		}

		assertHelper(t, testCase.CSVEntry, actual, testCaseName)
	}
}

func TestMakeUnsupported(t *testing.T) {
	t.Parallel()

	testCases := fakes.TransactionTestCasesUnsupported
	factory := transaction.NewCSVEntryFactory(log.New())

	for testCaseName, testCase := range testCases {
		_, err := factory.Make(testCase.Transaction)

		assert.Error(t, err, fmt.Sprintf("case '%s'", testCaseName))
	}
}

func TestMakeUnknown(t *testing.T) {
	t.Parallel()

	testCases := fakes.TransactionTestCasesUnknown
	factory := transaction.NewCSVEntryFactory(log.New())

	for testCaseName, testCase := range testCases {
		_, err := factory.Make(testCase.Transaction)

		assert.Error(t, err, fmt.Sprintf("case '%s'", testCaseName))
	}
}

// helper to assert float64 fields.
func assertHelper(t *testing.T, expected, actual filesystem.CSVEntry, testCaseName string) {
	t.Helper()

	assert.Equal(
		t,
		expected.Type,
		actual.Type,
		fmt.Sprintf("case '%s': type does not match", testCaseName),
	)

	assert.NotEqual(
		t,
		internal.DateTime{},
		actual.Timestamp,
		fmt.Sprintf("case '%s': timestamp is empty", testCaseName),
	)

	assert.Equal(
		t,
		expected.Timestamp,
		actual.Timestamp,
		fmt.Sprintf("case '%s': timestamp does not match", testCaseName),
	)

	assert.Equal(
		t,
		floatToStr(expected.Shares),
		floatToStr(actual.Shares),
		fmt.Sprintf("case '%s': shares amount does not match", testCaseName),
	)
	assert.Equal(
		t,
		floatToStr(expected.Rate),
		floatToStr(actual.Rate),
		fmt.Sprintf("case '%s': rate does not match", testCaseName),
	)
	assert.Equal(
		t,
		floatToStr(expected.Commission),
		floatToStr(actual.Commission),
		fmt.Sprintf("case '%s': commission amount does not match", testCaseName),
	)
	assert.Equal(
		t,
		floatToStr(expected.Yield),
		floatToStr(actual.Yield),
		fmt.Sprintf("case '%s': yield does not match", testCaseName),
	)
	assert.Equal(
		t,
		floatToStr(expected.Profit),
		floatToStr(actual.Profit),
		fmt.Sprintf("case '%s': profit does not match", testCaseName),
	)
	assert.Equal(
		t,
		floatToStr(expected.Debit),
		floatToStr(actual.Debit),
		fmt.Sprintf("case '%s': debit amount does not match", testCaseName),
	)
	assert.Equal(
		t,
		floatToStr(expected.Credit),
		floatToStr(actual.Credit),
		fmt.Sprintf("case '%s': credit amount does not match", testCaseName),
	)
}

// converts float64 to string to simplify assertions since data written into CSV will be in string anyway.
func floatToStr(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
