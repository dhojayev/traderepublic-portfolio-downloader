package internal_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
)

func TestParseTimestamp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input string
		year  int
	}{
		{
			input: "2023-12-13T12:44:28.857+0000",
			year:  2023,
		},
		{
			input: "2025-01-01T09:30:52.2408+01:00",
			year:  2025,
		},
		{
			input: "2024-11-01T05:46:41.70631+01:00",
			year:  2024,
		},
	}

	for i, testCase := range testCases {
		actual, err := internal.ParseTimestamp(testCase.input)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.year, actual.Year(), fmt.Sprintf("case %d", i))
	}
}
