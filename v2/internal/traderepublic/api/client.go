package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/restclient"
)

const (
	// HTTP status code threshold for error responses.
	statusCodeError = http.StatusBadRequest
)

// Client is a client that uses the generated OpenAPI client.
type Client struct {
	client *restclient.ClientWithResponses
}

// NewClient creates a new client that uses the generated OpenAPI client.
func NewClient() (*Client, error) {
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
	}, nil
}

// Login logs in with phone number and PIN.
func (c *Client) Login(requestBody restclient.APILoginRequest) (string, error) {
	// Make the login request
	resp, err := c.client.LoginWithResponse(context.Background(), requestBody)
	if err != nil {
		return "", fmt.Errorf("could not login: %w", err)
	}

	// Check for error response
	if resp.StatusCode() >= statusCodeError {
		return "", fmt.Errorf(
			"login failed with status code %d: %s",
			resp.StatusCode(),
			string(resp.Body),
		)
	}

	// Extract the process ID from the response
	var processID string
	if resp.JSON200 != nil && resp.JSON200.ProcessId != nil {
		processID = *resp.JSON200.ProcessId
	}

	return processID, nil
}

// PostOTP verifies the OTP.
func (c *Client) PostOTP(processID, otp string) ([]*http.Cookie, error) {
	if processID == "" {
		return nil, errors.New("processID cannot be empty")
	}

	// Make the OTP verification request
	resp, err := c.client.VerifyOTPWithResponse(context.Background(), processID, otp)
	if err != nil {
		return nil, fmt.Errorf("could not validate otp: %w", err)
	}

	// Check for error response
	if resp.StatusCode() >= statusCodeError {
		return nil, fmt.Errorf(
			"OTP verification failed with status code %d: %s",
			resp.StatusCode(),
			string(resp.Body),
		)
	}

	// Return all cookies from the response
	return resp.HTTPResponse.Cookies(), nil
}
