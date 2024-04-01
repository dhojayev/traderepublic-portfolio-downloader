package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenFromHeader(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		header        http.Header
		mustReturnErr bool
		expected      string
	}{
		{
			header: http.Header{
				"Content-type": []string{"application/json"},
			},
			mustReturnErr: true,
		},
		{
			header: http.Header{
				"Content-type": []string{"application/json"},
				"Set-Cookie": []string{
					"JSESSIONID=22192A210742959A362F7820ECC81311; Path=/; Secure; HttpOnly",
				},
			},
			mustReturnErr: true,
		},
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
	}

	assert := assert.New(t)

	for _, testCase := range testCases {
		token, err := NewTokenFromHeader("session", testCase.header)

		if !testCase.mustReturnErr {
			assert.Nil(err)
		}

		assert.Equal("session", token.Name())
		assert.Equal(testCase.expected, token.Value())
	}
}
