package transaction

import (
	"strconv"
	"testing"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFromPurchase(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		purchase Purchase
		expected filesystem.CSVEntry
	}{
		// purchased for 501 (including 1 eur commission)
		{
			purchase: NewPurchase(
				NewBaseTransaction("test-id", "test-status", time.Now()),
				NewAsset("test-instrument", "test-asset-name", 5.186721),
				NewMonetaryValues(96.40, 1, 501),
				[]Document{NewDocument("test-doc-id", "test-url", "test-date", "test-title")},
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

	assert := assert.New(t)
	factory := NewCSVEntryFactory(log.New())

	for _, testCase := range testCases {
		actual := factory.FromPurchase(testCase.purchase)

		assertFloat64Fields(assert, testCase.expected, actual)
	}
}

func TestFromSale(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		sale     Sale
		expected filesystem.CSVEntry
	}{
		// purchased for 258 (including 2 commissions of 1 eur), sold with profit.
		{
			sale: NewSale(
				43.9,
				113.25,
				NewPurchase(
					NewBaseTransaction("test-id", "test-status", time.Now()),
					NewAsset("test-instrument", "test-asset-name", 56.065306),
					NewMonetaryValues(6.62, 1, 370.25),
					[]Document{NewDocument("test-doc-id", "test-url", "test-date", "test-title")},
				),
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
			sale: NewSale(
				-0.62,
				-11.28,
				NewPurchase(
					NewBaseTransaction("test-id", "test-status", time.Now()),
					NewAsset("test-instrument", "test-asset-name", 21.272454),
					NewMonetaryValues(85.48, 1, 1817.27),
					[]Document{NewDocument("test-doc-id", "test-url", "test-date", "test-title")},
				),
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

	assert := assert.New(t)
	factory := NewCSVEntryFactory(log.New())

	for _, testCase := range testCases {
		actual := factory.FromSale(testCase.sale)

		assertFloat64Fields(assert, testCase.expected, actual)
	}
}

// helper to assert float64 fields.
func assertFloat64Fields(assert *assert.Assertions, expected, actual filesystem.CSVEntry) {
	assert.Equal(floatToStr(expected.Shares), floatToStr(actual.Shares), "shares amount does not match")
	assert.Equal(floatToStr(expected.Rate), floatToStr(actual.Rate), "rate does not match")
	assert.Equal(floatToStr(expected.Commission), floatToStr(actual.Commission), "commission amount does not match")
	assert.Equal(floatToStr(expected.Yield), floatToStr(actual.Yield), "yield does not match")
	assert.Equal(floatToStr(expected.Profit), floatToStr(actual.Profit), "profit does not match")
	assert.Equal(floatToStr(expected.Debit), floatToStr(actual.Debit), "debit amount does not match")
	assert.Equal(floatToStr(expected.Credit), floatToStr(actual.Credit), "credit amount does not match")
	assert.Equal(floatToStr(expected.InvestedAmount), floatToStr(actual.InvestedAmount), "invested amount does not match")
}

// converts float64 to string to simplify assertions since data written into CSV will be in string anyway.
func floatToStr(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
