package internal_test

import (
	"fmt"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshallCSV(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"25 Sep 23 08:45 +0000",
		"18 Dec 23 12:23 +0100",
		"01 Oct 24 05:30 +0200",
	}

	for i, input := range testCases {
		datetime := internal.DateTime{}
		err := datetime.UnmarshalCSV(input)

		assert.NoError(t, err, fmt.Sprintf("Failed with test case index %d", i))
	}
}
