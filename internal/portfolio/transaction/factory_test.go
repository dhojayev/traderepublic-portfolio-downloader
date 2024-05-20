package transaction_test

import (
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

func TestFromPurchase(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		trn      transaction.Model
		expected filesystem.CSVEntry
	}{
		// purchased for 501 (including 1 eur commission)
		{
			trn: transaction.NewTransaction(
				"test-id", transaction.TypePurchase, "test-status", 0, 0, 5.186721, 96.40, 1, 501, time.Now(),
				transaction.NewInstrument("test-instrument", "test-asset-name"),
				[]document.Model{document.NewModel("test-doc-id", "test-url", "test-date", "test-title")},
			),
			expected: filesystem.NewCSVEntry(
				"test-id",
				"test-status",
				"Purchase",
				"Other",
				"test-asset-name",
				"test-instrument",
				5.186721,
				96.40,
				0,
				0,
				1,
				501,
				0,
				500,
				internal.DateTime{Time: time.Now()},
			),
		},
	}

	factory := transaction.NewCSVEntryFactory(log.New())

	for _, testCase := range testCases {
		actual, err := factory.Make(testCase.trn)

		assert.Nil(t, err)
		assertFloat64Fields(t, testCase.expected, actual)
	}
}

func TestFromSale(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		trn      transaction.Model
		expected filesystem.CSVEntry
	}{
		// purchased for 258 (including 2 commissions of 1 eur), sold with profit.
		{
			trn: transaction.NewTransaction(
				"test-id", transaction.TypeSale, "test-status", 43.9, 113.25, 56.065306, 6.62, 1, 370.25, time.Now(),
				transaction.NewInstrument("test-instrument", "test-asset-name"),
				[]document.Model{document.NewModel("test-doc-id", "test-url", "test-date", "test-title")},
			),
			expected: filesystem.NewCSVEntry(
				"test-id",
				"test-status",
				"Purchase",
				"Other",
				"test-asset-name",
				"test-instrument",
				-56.065306,
				6.62,
				43.9,
				113.25,
				1,
				0,
				370.25,
				-258,
				internal.DateTime{Time: time.Now()},
			),
		},
		// purchased for 1829.55 (including 5 commissions of 1 eur), sold with loss.
		{
			trn: transaction.NewTransaction(
				"test-id", transaction.TypeSale, "test-status", -0.62, -11.28, 21.272454, 85.48, 1, 1817.27, time.Now(),
				transaction.NewInstrument("test-instrument", "test-asset-name"),
				[]document.Model{document.NewModel("test-doc-id", "test-url", "test-date", "test-title")},
			),
			expected: filesystem.NewCSVEntry(
				"test-id",
				"test-status",
				"Purchase",
				"Other",
				"test-asset-name",
				"test-instrument",
				-21.272454,
				85.48,
				-0.62,
				-11.28,
				1,
				0,
				1817.27,
				-1829.55,
				internal.DateTime{Time: time.Now()},
			),
		},
	}

	factory := transaction.NewCSVEntryFactory(log.New())

	for _, testCase := range testCases {
		actual, err := factory.Make(testCase.trn)

		assert.Nil(t, err)
		assertFloat64Fields(t, testCase.expected, actual)
	}
}

// helper to assert float64 fields.
func assertFloat64Fields(t *testing.T, expected, actual filesystem.CSVEntry) {
	t.Helper()

	assert.Equal(t, floatToStr(expected.Shares), floatToStr(actual.Shares), "shares amount does not match")
	assert.Equal(t, floatToStr(expected.Rate), floatToStr(actual.Rate), "rate does not match")
	assert.Equal(t, floatToStr(expected.Commission), floatToStr(actual.Commission), "commission amount does not match")
	assert.Equal(t, floatToStr(expected.Yield), floatToStr(actual.Yield), "yield does not match")
	assert.Equal(t, floatToStr(expected.Profit), floatToStr(actual.Profit), "profit does not match")
	assert.Equal(t, floatToStr(expected.Debit), floatToStr(actual.Debit), "debit amount does not match")
	assert.Equal(t, floatToStr(expected.Credit), floatToStr(actual.Credit), "credit amount does not match")
	assert.Equal(t, floatToStr(expected.InvestedAmount), floatToStr(actual.InvestedAmount), "invested amount does not match")
}

// converts float64 to string to simplify assertions since data written into CSV will be in string anyway.
func floatToStr(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
