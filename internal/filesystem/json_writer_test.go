package filesystem_test

import (
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFilename(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		dir      string
		dataMap  map[string]any
		expected string
		hasErr   bool
	}{
		{
			dir:      "purchase",
			dataMap:  map[string]any{"id": "test-id"},
			expected: "test-id",
		},
		{
			dir:      "purchase",
			dataMap:  map[string]any{"id": 1},
			expected: "",
			hasErr:   true,
		},
		{
			dir:      "transactions",
			dataMap:  map[string]any{"test-field": "test-val"},
			expected: "page-1",
		},
		{
			dir:      "sales",
			dataMap:  map[string]any{"test-field": "test-val"},
			expected: "page-1",
		},
		{
			dir:      "transactions",
			dataMap:  map[string]any{"test-field": "test-val"},
			expected: "page-2",
		},
		{
			dir:      "sales",
			dataMap:  map[string]any{"test-field": "test-val"},
			expected: "page-2",
		},
	}

	logger := log.New()
	writer := filesystem.NewJSONWriter(logger)
	asrt := assert.New(t)

	for _, testCase := range testCases {
		actual, err := writer.GenerateFilename(testCase.dir, testCase.dataMap)

		asrt.Equal(testCase.expected, actual)

		if testCase.hasErr {
			asrt.NotNil(err)
		}
	}
}
