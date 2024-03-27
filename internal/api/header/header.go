package header

import (
	"net/http"
	"strings"
)

type Headers http.Header

func NewHeaders() Headers {
	return Headers(http.Header{
		"User-Agent": {strings.Join([]string{
			"Mozilla/5.0",
			"(Macintosh; Intel Mac OS X 10_15_7)",
			"AppleWebKit/605.1.15 (KHTML, like Gecko)",
			"Version/17.4 Safari/605.1.15",
		}, " ")},
	})
}

func (h Headers) With(key, value string) Headers {
	if _, found := h[key]; found {
		h[key] = append(h[key], value)

		return h
	}

	h[key] = []string{value}

	return h
}

func (h Headers) WithContentTypeJSON() Headers {
	h["Content-Type"] = []string{"application/json"}

	return h
}

func (h Headers) WithRefreshToken(token string) Headers {
	h["Cookie"] = []string{"tr_refresh=" + token}

	return h
}

func (h Headers) AsHTTPHeader() http.Header {
	return http.Header(h)
}
