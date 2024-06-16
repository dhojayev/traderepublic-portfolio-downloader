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

	for i, testCase := range testCases {
		actual, err := factory.Make(testCase.Transaction)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		if err != nil {
			continue
		}

		assertHelper(t, testCase.CSVEntry, actual, i)
	}
}

func TestMakeUnsupported(t *testing.T) {
	t.Parallel()

	testCases := fakes.TransactionTestCasesUnsupported
	factory := transaction.NewCSVEntryFactory(log.New())

	for i, testCase := range testCases {
		_, err := factory.Make(testCase.Transaction)

		assert.Error(t, err, fmt.Sprintf("case %d", i))
	}
}

func TestMakeUnknown(t *testing.T) {
	t.Parallel()

	testCases := fakes.TransactionTestCasesUnknown
	factory := transaction.NewCSVEntryFactory(log.New())

	for i, testCase := range testCases {
		_, err := factory.Make(testCase.Transaction)

		assert.Error(t, err, fmt.Sprintf("case %d", i))
	}
}

// helper to assert float64 fields.
func assertHelper(t *testing.T, expected, actual filesystem.CSVEntry, testCaseNum int) {
	t.Helper()

	assert.Equal(
		t,
		expected.Type,
		actual.Type,
		fmt.Sprintf("case %d: type does not match", testCaseNum),
	)

	assert.NotEqual(
		t,
		internal.DateTime{},
		actual.Timestamp,
		fmt.Sprintf("case %d: timestamp is empty", testCaseNum),
	)

	assert.Equal(
		t,
		expected.Timestamp,
		actual.Timestamp,
		fmt.Sprintf("case %d: timestamp does not match", testCaseNum),
	)

	assert.Equal(
		t,
		floatToStr(expected.Shares),
		floatToStr(actual.Shares),
		fmt.Sprintf("case %d: shares amount does not match", testCaseNum),
	)
	assert.Equal(
		t,
		floatToStr(expected.Rate),
		floatToStr(actual.Rate),
		fmt.Sprintf("case %d: rate does not match", testCaseNum),
	)
	assert.Equal(
		t,
		floatToStr(expected.Commission),
		floatToStr(actual.Commission),
		fmt.Sprintf("case %d: commission amount does not match", testCaseNum),
	)
	assert.Equal(
		t,
		floatToStr(expected.Yield),
		floatToStr(actual.Yield),
		fmt.Sprintf("case %d: yield does not match", testCaseNum),
	)
	assert.Equal(
		t,
		floatToStr(expected.Profit),
		floatToStr(actual.Profit),
		fmt.Sprintf("case %d: profit does not match", testCaseNum),
	)
	assert.Equal(
		t,
		floatToStr(expected.Debit),
		floatToStr(actual.Debit),
		fmt.Sprintf("case %d: debit amount does not match", testCaseNum),
	)
	assert.Equal(
		t,
		floatToStr(expected.Credit),
		floatToStr(actual.Credit),
		fmt.Sprintf("case %d: credit amount does not match", testCaseNum),
	)
}

// converts float64 to string to simplify assertions since data written into CSV will be in string anyway.
func floatToStr(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
