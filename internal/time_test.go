package internal_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
)

func TestParseTimestamp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input      string
		year       int
		month      time.Month
		day        int
		hour       int
		minute     int
		second     int
		nanosecond int
	}{
		{
			input:      "2023-12-13T12:44:28.857+0000",
			year:       2023,
			month:      12,
			day:        13,
			hour:       12,
			minute:     44,
			second:     28,
			nanosecond: 857000000,
		},
		{
			input:      "2025-01-01T09:30:52.2408+01:00",
			year:       2025,
			month:      1,
			day:        1,
			hour:       9,
			minute:     30,
			second:     52,
			nanosecond: 240800000,
		},
		{
			input:      "2024-11-01T05:46:41.70631+01:00",
			year:       2024,
			month:      11,
			day:        1,
			hour:       5,
			minute:     46,
			second:     41,
			nanosecond: 706310000,
		},
	}

	for i, testCase := range testCases {
		actual, err := internal.ParseTimestamp(testCase.input)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.year, actual.Year(), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.month, actual.Month(), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.day, actual.Day(), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.hour, actual.Hour(), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.minute, actual.Minute(), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.second, actual.Second(), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.nanosecond, actual.Nanosecond(), fmt.Sprintf("case %d", i))
	}
}
