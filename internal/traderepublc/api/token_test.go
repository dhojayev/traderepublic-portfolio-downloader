package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api"
)

func TestItCanCreateNewTokenFromHeader(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		header   http.Header
		expected string
	}{
		{
			header: http.Header{
				"Content-type": []string{"application/json"},
				"Set-Cookie": []string{
					"JSESSIONID=22192A210742959A362F7820ECC81311; Path=/; Secure; HttpOnly",
					"tr_session=qPzbun8CCMGhaWZmFySyziGgmdDn97uC0pWcqDJNWcM8eUwF5mq5UJmMzjvdEwDj; Path=/; Secure; HttpOnly; SameSite=Strict",
				},
			},
			expected: "qPzbun8CCMGhaWZmFySyziGgmdDn97uC0pWcqDJNWcM8eUwF5mq5UJmMzjvdEwDj",
		},
		{
			header: http.Header{
				"Content-type": []string{"application/json"},
				"Set-Cookie": []string{
					"JSESSIONID=22192A210742959A362F7820ECC81311; Path=/; Secure; HttpOnly",
					"tr_session=sOLtnIhoYnTF4AYDvQyG.ckTkgq5eQBKHz1dXET8PAjtdrtr4uGV2h815hViVe1fjPHNFC0SrXRgJbeyncKirTPfJfEFyMyqDrxmJQdd6tewLBuARVUCKQveEeYLSmMm3Yk2SVSpQdxCzrCiepByqLV8C9GV6gd4gxfR2P8PQ99TYXGZ5mpUugdQEuBrDa2x3fArzBjGWiukkexaG70SECAHXYgX7DnW0NT4bAzr32E258ycRXwDm2fGLdKeHLFvZktuT7MXK8rX4M5FJaFYQv6tQEH5DV6HQvgpcacGMgGqvLygR5LSfuwzD9MUmhkU54daCRxABegfkDJ9mzBhfg5Z7CTd7NdbgJRCPAnHdFdEdQZCDyt31ALQNuatXbTUk2gy2XVC0ZEKYaAUWLF49E2xaGMyyW3tQ3V9KwFTpPf.QdjFGjjZaf5PTMnKPtG51DFUwctzpfCW8LLT0rxEzzZbrVJtZJbzxUQWiYfJPa5EXWz0Sz7Hu9NmW5jnUREfAw; Path=/; Secure; HttpOnly; SameSite=Strict",
				},
			},
			expected: "sOLtnIhoYnTF4AYDvQyG.ckTkgq5eQBKHz1dXET8PAjtdrtr4uGV2h815hViVe1fjPHNFC0SrXRgJbeyncKirTPfJfEFyMyqDrxmJQdd6tewLBuARVUCKQveEeYLSmMm3Yk2SVSpQdxCzrCiepByqLV8C9GV6gd4gxfR2P8PQ99TYXGZ5mpUugdQEuBrDa2x3fArzBjGWiukkexaG70SECAHXYgX7DnW0NT4bAzr32E258ycRXwDm2fGLdKeHLFvZktuT7MXK8rX4M5FJaFYQv6tQEH5DV6HQvgpcacGMgGqvLygR5LSfuwzD9MUmhkU54daCRxABegfkDJ9mzBhfg5Z7CTd7NdbgJRCPAnHdFdEdQZCDyt31ALQNuatXbTUk2gy2XVC0ZEKYaAUWLF49E2xaGMyyW3tQ3V9KwFTpPf.QdjFGjjZaf5PTMnKPtG51DFUwctzpfCW8LLT0rxEzzZbrVJtZJbzxUQWiYfJPa5EXWz0Sz7Hu9NmW5jnUREfAw",
		},
		{
			header: http.Header{
				"Set-Cookie": []string{
					"JSESSIONID=7AA9CC2A4DDF395AAD5920F097C7F559; Path=/; Secure; HttpOnly",
					"tr_session=eyJhbGciOiJFUzI1NiJ9.eyJqdXJpc2RpY3Rpb24iOiJERSIsInNlc3Npb25JZCI6ImQxNTM3Y2U2LWYyZGQtNDhlZC1hN2IzLWY2OTdlMTNlMTAyNyIsInR5cGUiOiJTRVNTSU9OIiwiZmVhdHVyZXNFbmFibGVkIjpbeyJmZWF0dXJlIjoiY3VycmVudEFjY291bnRBdmFpbGFibGUifSx7ImZlYXR1cmUiOiJ0ckliYW5BY3RpdmF0ZWQifSx7ImZlYXR1cmUiOiJjcnlwdG8ifSx7ImZlYXR1cmUiOiJjYXJkIn0seyJmZWF0dXJlIjoiYm9uZFRyYWRpbmcifSx7ImZlYXR1cmUiOiJwdkVuYWJsZWRGb3JDYXJkcyJ9XSwic3RhdHVzIjoiQUNUSVZFIiwiaWF0IjoxNzQyODE0MjYwLCJleHAiOjE3NDI4MTQ1NjAsInN1YiI6IjEzOGM1NjJiLWI2ODgtNDhkMi05OGU3LTRjMGI1NTZmOThiNyJ9.mkdoWAiTiyMP5IbJyrEY_ZzU--dQpUs-T0plXUpUGG_1wWU48BJLCK7_yLOc0j38h66z0cL7VJCtNOujlLGRww; Path=/; Secure; HttpOnly; SameSite=Strict",
				},
			},
			expected: "eyJhbGciOiJFUzI1NiJ9.eyJqdXJpc2RpY3Rpb24iOiJERSIsInNlc3Npb25JZCI6ImQxNTM3Y2U2LWYyZGQtNDhlZC1hN2IzLWY2OTdlMTNlMTAyNyIsInR5cGUiOiJTRVNTSU9OIiwiZmVhdHVyZXNFbmFibGVkIjpbeyJmZWF0dXJlIjoiY3VycmVudEFjY291bnRBdmFpbGFibGUifSx7ImZlYXR1cmUiOiJ0ckliYW5BY3RpdmF0ZWQifSx7ImZlYXR1cmUiOiJjcnlwdG8ifSx7ImZlYXR1cmUiOiJjYXJkIn0seyJmZWF0dXJlIjoiYm9uZFRyYWRpbmcifSx7ImZlYXR1cmUiOiJwdkVuYWJsZWRGb3JDYXJkcyJ9XSwic3RhdHVzIjoiQUNUSVZFIiwiaWF0IjoxNzQyODE0MjYwLCJleHAiOjE3NDI4MTQ1NjAsInN1YiI6IjEzOGM1NjJiLWI2ODgtNDhkMi05OGU3LTRjMGI1NTZmOThiNyJ9.mkdoWAiTiyMP5IbJyrEY_ZzU--dQpUs-T0plXUpUGG_1wWU48BJLCK7_yLOc0j38h66z0cL7VJCtNOujlLGRww",
		},
	}

	for i, testCase := range testCases {
		token, err := api.NewTokenFromHeader("session", testCase.header)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, "session", token.Name(), fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.expected, token.Value(), fmt.Sprintf("case %d", i))
	}
}

func TestItReturnsErrorOnNoSessionInHeader(t *testing.T) {
	t.Parallel()

	testCases := []http.Header{
		{
			"Content-type": []string{"application/json"},
		},
		{
			"Content-type": []string{"application/json"},
			"Set-Cookie": []string{
				"JSESSIONID=22192A210742959A362F7820ECC81311; Path=/; Secure; HttpOnly",
			},
		},
	}

	for i, testCase := range testCases {
		_, err := api.NewTokenFromHeader("session", testCase)

		assert.Error(t, err, fmt.Sprintf("case %d", i))
	}
}
