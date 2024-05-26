package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/header"
)

const (
	baseURL = "https://api.traderepublic.com/api/v1"
)

type Client struct {
	logger *log.Logger
}

func NewClient(logger *log.Logger) Client {
	return Client{
		logger: logger,
	}
}

func (c *Client) Login(requestBody LoginRequest, refreshToken Token) (LoginResponse, Token, error) {
	var responseBody LoginResponse

	sessionToken := NewToken(TokenNameSession, "")

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return responseBody, sessionToken, fmt.Errorf("could not marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		fmt.Sprintf("%s/%s", baseURL, "auth/web/login"),
		bytes.NewReader(requestBodyBytes),
	)
	if err != nil {
		return responseBody, sessionToken, fmt.Errorf("could not create login request: %w", err)
	}

	h := header.NewHeaders().WithContentTypeJSON()
	if refreshToken.Value() != "" {
		h = h.WithRefreshToken(refreshToken.Value())
	}

	req.Header = h.AsHTTPHeader()

	resp, err := c.request(req)
	if err != nil {
		return responseBody, sessionToken, err
	}

	defer func() { _ = resp.Body.Close() }()

	sessionToken, _ = NewTokenFromHeader(TokenNameSession, resp.Header)

	if err = c.readResponseBody(resp.Body, &responseBody); err != nil {
		return responseBody, sessionToken, err
	}

	return responseBody, sessionToken, nil
}

func (c *Client) PostOTP(processID, otp string) (Token, Token, error) {
	sessionToken := NewToken(TokenNameSession, "")
	refreshToken := NewToken(TokenNameRefresh, "")

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		fmt.Sprintf("%s/%s/%s/%s", baseURL, "auth/web/login", processID, otp),
		nil,
	)
	if err != nil {
		return sessionToken, refreshToken, fmt.Errorf("could not create otp request: %w", err)
	}

	req.Header = header.NewHeaders().WithContentTypeJSON().AsHTTPHeader()

	resp, err := c.request(req)
	if err != nil {
		return sessionToken, refreshToken, err
	}

	defer func() { _ = resp.Body.Close() }()

	c.logger.Tracef("%#v", resp.Header)

	sessionToken, err = NewTokenFromHeader(TokenNameSession, resp.Header)
	if err != nil {
		return sessionToken, refreshToken, fmt.Errorf("could not parse session token from header: %w", err)
	}

	refreshToken, err = NewTokenFromHeader(TokenNameRefresh, resp.Header)
	if err != nil {
		return sessionToken, refreshToken, fmt.Errorf("could not parse refresh token from header: %w", err)
	}

	c.logger.Debug("received session and refresh tokens")

	return sessionToken, refreshToken, nil
}

func (c *Client) Session(refreshToken Token) (Token, error) {
	sessionToken := NewToken(TokenNameSession, "")
	url := fmt.Sprintf("%s/%s", baseURL, "auth/web/session")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return sessionToken, fmt.Errorf("could not create session request: %w", err)
	}

	req.Header = header.NewHeaders().WithContentTypeJSON().WithRefreshToken(refreshToken.Value()).AsHTTPHeader()

	resp, err := c.request(req)
	if err != nil {
		return sessionToken, fmt.Errorf("could not request session endpoint: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	sessionToken, _ = NewTokenFromHeader(TokenNameSession, resp.Header)

	return sessionToken, nil
}

func (c *Client) request(req *http.Request) (*http.Response, error) {
	c.logger.
		WithFields(log.Fields{
			"method": req.Method,
			"url":    req.URL.String(),
		}).
		Trace("executing request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("could not make request: %w", err)
	}

	if resp.StatusCode < http.StatusBadRequest {
		return resp, nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("could not parse response body: %w", err)
	}

	return resp, fmt.Errorf("request failed with status code '%d': %s", resp.StatusCode, respBody)
}

func (c *Client) readResponseBody(body io.ReadCloser, v any) error {
	responseBodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	c.logger.
		WithField("body", string(responseBodyBytes)).
		Debug("received success response")

	if err = json.Unmarshal(responseBodyBytes, v); err != nil {
		return fmt.Errorf("could not unmarshal login response: %w", err)
	}

	return nil
}
