package timelinedetails_test

import (
	"fmt"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/timelinedetails"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFloatFromResponse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "Du hast 500,00 € per Lastschrift hinzugefügt", expected: 500.00},
		{input: "Du hast 0,123989123 € per Lastschrift hinzugefügt", expected: 0.123989123},
		{input: "Du hast 200,00 € erhalten", expected: 200.00},
		{input: "Du hast 1,00 € erhalten", expected: 1.00},
		{input: "Du hast 280,85 €  erhalten", expected: 280.85},
		{input: "Du hast 66,60 EUR erhalten", expected: 66.60},
		{input: "Du hast 1.000,00 € erhalten", expected: 1000.00},
		{input: "Du hast 1.921,89 €  investiert", expected: 1921.89},
		{input: "Du hast 10.000,00 € erhalten", expected: 10000.00},
		{input: "500,00 €", expected: 500},
		{input: "0,12382889 €", expected: 0.12382889},
		{input: "200,00 €", expected: 200},
		{input: "1,00 €", expected: 1},
		{input: "280,85 € ", expected: 280.85},
		{input: "66,60 EUR", expected: 66.6},
		{input: "1.000,00 €", expected: 1000},
		{input: "1.921,89 € ", expected: 1921.89},
		{input: "10.000,00 €", expected: 10000},
		{input: "9 %", expected: 9},
		{input: "138.26 €", expected: 138.26},
		{input: "500.00", expected: 500},
		{input: "0.0234898", expected: 0.0234898},
		{input: "200.00", expected: 200},
		{input: "1.00", expected: 1},
		{input: "280.85", expected: 280.85},
		{input: "66.60", expected: 66.6},
		{input: "1000.00", expected: 1000},
		{input: "1921.89", expected: 1921.89},
		{input: "10000.00", expected: 10000},
		{input: "138.26 €", expected: 138.26},
	}

	for i, testCase := range testCases {
		actual, err := timelinedetails.ParseFloatFromResponse(testCase.input)

		require.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.expected, actual, fmt.Sprintf("case %d", i))
	}
}
