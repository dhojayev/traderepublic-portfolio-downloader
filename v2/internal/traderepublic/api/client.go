package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/restclient"
)

const (
	// HTTP status code threshold for error responses.
	statusCodeError = http.StatusBadRequest

	// Cookie names.
	cookieNameSession = "tr_session"
	cookieNameRefresh = "tr_refresh"
)

// Client is a client that uses the generated OpenAPI client.
type Client struct {
	client *restclient.ClientWithResponses
	logger *slog.Logger
}

// NewClient creates a new client that uses the generated OpenAPI client.
func NewClient(logger *slog.Logger) (*Client, error) {
	// Create a request editor to add common headers
	reqEditor := func(_ context.Context, req *http.Request) error {
		req.Header.Set("User-Agent", internal.HTTPUserAgent)
		req.Header.Set("Content-Type", "application/json")

		return nil
	}

	// Create the client with the base URL and request editor
	client, err := restclient.NewClientWithResponses(
		internal.RestAPIBaseURI,
		restclient.WithRequestEditorFn(reqEditor),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create REST client: %w", err)
	}

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

// Login logs in with phone number and PIN.
func (c *Client) Login(
	requestBody restclient.APILoginRequest,
	refreshToken Token,
) (restclient.APILoginResponse, Token, error) {
	var responseBody restclient.APILoginResponse

	sessionToken := NewToken(TokenNameSession, "")

	// Create a request editor to add the refresh token if available
	reqEditor := func(_ context.Context, req *http.Request) error {
		if refreshToken.Value() != "" {
			req.Header.Set("Cookie", cookieNameRefresh+"="+refreshToken.Value())
		}

		return nil
	}

	// Make the login request
	resp, err := c.client.LoginWithResponse(context.Background(), requestBody, reqEditor)
	if err != nil {
		return responseBody, sessionToken, fmt.Errorf("could not login: %w", err)
	}

	// Check for error response
	if resp.StatusCode() >= statusCodeError {
		return responseBody, sessionToken, fmt.Errorf(
			"login failed with status code %d: %s",
			resp.StatusCode(),
			string(resp.Body),
		)
	}

	// Extract the session token from the response
	for _, cookie := range resp.HTTPResponse.Cookies() {
		if cookie.Name == cookieNameSession {
			sessionToken = NewToken(TokenNameSession, cookie.Value)

			break
		}
	}

	// Set the response body
	if resp.JSON200 != nil {
		responseBody = *resp.JSON200
	}

	return responseBody, sessionToken, nil
}

// PostOTP verifies the OTP.
func (c *Client) PostOTP(processID, otp string) (Token, Token, error) {
	sessionToken := NewToken(TokenNameSession, "")
	refreshToken := NewToken(TokenNameRefresh, "")

	if processID == "" {
		return sessionToken, refreshToken, errors.New("processID cannot be empty")
	}

	// Make the OTP verification request
	resp, err := c.client.VerifyOTPWithResponse(context.Background(), processID, otp)
	if err != nil {
		return sessionToken, refreshToken, fmt.Errorf("could not validate otp: %w", err)
	}

	// Check for error response
	if resp.StatusCode() >= statusCodeError {
		return sessionToken, refreshToken, fmt.Errorf(
			"OTP verification failed with status code %d: %s",
			resp.StatusCode(),
			string(resp.Body),
		)
	}

	// Extract the tokens from the response
	for _, cookie := range resp.HTTPResponse.Cookies() {
		switch cookie.Name {
		case cookieNameSession:
			sessionToken = NewToken(TokenNameSession, cookie.Value)
		case cookieNameRefresh:
			refreshToken = NewToken(TokenNameRefresh, cookie.Value)
		}
	}

	return sessionToken, refreshToken, nil
}

// Session refreshes the session token.
func (c *Client) Session(refreshToken Token) (Token, error) {
	sessionToken := NewToken(TokenNameSession, "")

	// Create a request editor to add the refresh token
	reqEditor := func(_ context.Context, req *http.Request) error {
		req.Header.Set("Cookie", "tr_refresh="+refreshToken.Value())

		return nil
	}

	// Make the session refresh request
	resp, err := c.client.RefreshSessionWithResponse(context.Background(), reqEditor)
	if err != nil {
		return sessionToken, fmt.Errorf("could not refresh session: %w", err)
	}

	// Check for error response
	if resp.StatusCode() >= statusCodeError {
		return sessionToken, fmt.Errorf(
			"session refresh failed with status code %d: %s",
			resp.StatusCode(),
			string(resp.Body),
		)
	}

	// Extract the session token from the response
	for _, cookie := range resp.HTTPResponse.Cookies() {
		if cookie.Name == cookieNameSession {
			sessionToken = NewToken(TokenNameSession, cookie.Value)

			break
		}
	}

	return sessionToken, nil
}
