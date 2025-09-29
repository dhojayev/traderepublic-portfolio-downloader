//go:generate go tool mockgen -source=client.go -destination client_mock.go -package=auth

package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/restclient"
)

type Client struct {
	inputHandler console.InputHandlerInterface
	apiClient    api.ClientInterface
}

func NewClient(inputHandler console.InputHandlerInterface, apiClient api.ClientInterface) *Client {
	return &Client{
		inputHandler: inputHandler,
		apiClient:    apiClient,
	}
}

func (c *Client) Login() (Token, error) {
	var token Token

	processID, err := c.ObtainProcessID()
	if err != nil {
		return token, fmt.Errorf("could not obtain process ID: %w", err)
	}

	token, err = c.ProvideOTP(processID)
	if err != nil {
		return token, fmt.Errorf("could not provide OTP: %w", err)
	}

	return token, nil
}

func (c *Client) ObtainProcessID() (ProcessID, error) {
	phoneNumber, err := c.inputHandler.GetPhoneNumber()
	if err != nil {
		return "", fmt.Errorf("failed to get phone number: %w", err)
	}

	pin, err := c.inputHandler.GetPIN()
	if err != nil {
		return "", fmt.Errorf("failed to get PIN: %w", err)
	}

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

func (c *Client) ProvideOTP(processID ProcessID) (Token, error) {
	var token Token

	if processID == "" {
		return token, errors.New("processID cannot be empty")
	}

	otp, err := c.inputHandler.GetOTP()
	if err != nil {
		return token, fmt.Errorf("failed to get OTP: %w", err)
	}

	// Call the API client's PostOTP method
	cookies, err := c.apiClient.PostOTP(string(processID), string(otp))
	if err != nil {
		return token, fmt.Errorf("could not validate otp: %w", err)
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
