package filesystem_test

import (
	"fmt"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
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
	logger.Out = io.Discard

	writer := filesystem.NewJSONWriter(logger)

	for i, testCase := range testCases {
		actual, err := writer.GenerateFilename(testCase.dir, testCase.dataMap)
		assert.Equal(t, testCase.expected, actual, fmt.Sprintf("case %d", i))

		if testCase.hasErr {
			assert.NotNil(t, err, fmt.Sprintf("case %d", i))
		}
	}
}
