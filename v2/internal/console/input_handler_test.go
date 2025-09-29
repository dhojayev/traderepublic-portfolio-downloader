package console_test

import (
	"io"
	"os"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/console"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStdin replaces os.Stdin with a pipe that we can write to for testing.
func mockStdin(t *testing.T, input string) func() {
	t.Helper()

	// Save original stdin
	oldStdin := os.Stdin

	// Create a pipe
	read, write, err := os.Pipe()
	require.NoError(t, err, "Creating pipe should not error")

	// Replace stdin with our pipe
	os.Stdin = read

	// Write the input to the pipe
	_, err = write.Write([]byte(input))
	require.NoError(t, err, "Writing to pipe should not error")

	_ = write.Close()

	// Return a cleanup function
	return func() {
		os.Stdin = oldStdin
	}
}

func TestGetPhoneNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expected      string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid phone number",
			input:       "+491234567890\n",
			expected:    "+491234567890",
			expectError: false,
		},
		{
			name:          "Empty input",
			input:         "\n",
			expectError:   true,
			errorContains: "could not read phone number",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Capture stdout to prevent it from cluttering the test output
			oldStdout := os.Stdout
			readPipe, writePipe, _ := os.Pipe()
			os.Stdout = writePipe

			// Mock stdin
			cleanup := mockStdin(t, testCase.input)
			defer cleanup()

			// Create handler and call method
			handler := console.NewInputHandler()
			result, err := handler.GetPhoneNumber()

			// Restore stdout
			_ = writePipe.Close()
			os.Stdout = oldStdout
			_, _ = io.Copy(io.Discard, readPipe)

			// Check results
			if testCase.expectError {
				assert.Error(t, err, "Should return an error")

				if err != nil {
					assert.Contains(t, err.Error(), testCase.errorContains, "Error message should contain expected text")
				}

				return
			}

			require.NoError(t, err, "Should not return an error")
			assert.Equal(t, testCase.expected, result, "Result should match expected value")
		})
	}
}

func TestGetPIN(t *testing.T) {
	t.Parallel()
	// Skip this test for now as it requires terminal input
	t.Skip("Skipping test that requires terminal input")
}

func TestGetOTP(t *testing.T) {
	t.Parallel()
	// Skip this test for now as it requires terminal input
	t.Skip("Skipping test that requires terminal input")
}
