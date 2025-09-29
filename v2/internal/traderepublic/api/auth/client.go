//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client.go -destination client_mock.go -package=auth

package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/restclient"
)

type Client struct {
	apiClient api.ClientInterface
}

func NewClient(apiClient api.ClientInterface) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

func (c *Client) Login(phoneNumber PhoneNumber, pin Pin) (ProcessID, error) {
	// Create the login request
	request := restclient.APILoginRequest{
		PhoneNumber: string(phoneNumber),
		Pin:         string(pin),
	}

	// Call the API client's Login method
	processID, err := c.apiClient.Login(request)
	if err != nil {
		return "", fmt.Errorf("could not login: %w", err)
	}

	return ProcessID(processID), nil
}

func (c *Client) ProvideOTP(processID ProcessID, otp OTP) (Token, error) {
	if processID == "" {
		return Token{}, errors.New("processID cannot be empty")
	}

	// Call the API client's PostOTP method
	cookies, err := c.apiClient.PostOTP(string(processID), string(otp))
	if err != nil {
		return Token{}, fmt.Errorf("could not validate otp: %w", err)
	}

	// Extract session and refresh tokens from cookies and return them
	return ExtractTokenFromCookies(cookies), nil
}

// ExtractTokenFromCookies creates a Token from HTTP cookies.
func ExtractTokenFromCookies(cookies []*http.Cookie) Token {
	var sessionValue, refreshValue string

	for _, cookie := range cookies {
		switch cookie.Name {
		case "tr_session":
			sessionValue = cookie.Value
		case "tr_refresh":
			refreshValue = cookie.Value
		}
	}

	return NewTokenWithValues(sessionValue, refreshValue)
}
