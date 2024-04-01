//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client.go -destination client_mock.go -package=auth

package auth

import (
	"errors"
	"fmt"
	"io/fs"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
)

const (
	tickerTimeoutSec = 60
)

type (
	PhoneNumber string
	Pin         string
)

type ClientInterface interface {
	Login() (api.LoginResponse, error)
	ProvideOTP(processID, otp string) error
	SessionToken() api.Token
}

type Client struct {
	apiClient    api.Client
	logger       *log.Logger
	phoneNumber  PhoneNumber
	pin          Pin
	sessionToken api.Token
	refreshToken api.Token
}

func NewClient(phoneNumber PhoneNumber, pin Pin, apiClient api.Client, logger *log.Logger) (*Client, error) {
	client := &Client{
		phoneNumber: phoneNumber,
		pin:         pin,
		apiClient:   apiClient,
		logger:      logger,
	}

	sessionToken, err := api.NewTokenFromFile(api.TokenNameSession)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("could not read session token file: %w", err)
	}

	client.sessionToken = sessionToken

	refreshToken, err := api.NewTokenFromFile(api.TokenNameRefresh)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("could not read refresh token file: %w", err)
	}

	client.refreshToken = refreshToken

	ticker := time.NewTicker(tickerTimeoutSec * time.Second)

	go func() {
		for range ticker.C {
			if err = client.refreshSession(); err != nil {
				client.logger.Error(err)
			}
		}
	}()

	return client, nil
}

func (c *Client) Login() (api.LoginResponse, error) {
	resp, sessionToken, err := c.apiClient.Login(
		api.LoginRequest{
			PhoneNumber: string(c.phoneNumber),
			Pin:         string(c.pin),
		},
		c.refreshToken,
	)
	if err != nil {
		return resp, fmt.Errorf("could not login: %w", err)
	}

	if sessionToken.Value() != "" {
		c.sessionToken = sessionToken
	}

	return resp, nil
}

func (c *Client) ProvideOTP(processID, otp string) error {
	if processID == "" {
		return errors.New("processID cannot be empty")
	}

	sessionToken, refreshToken, err := c.apiClient.PostOTP(processID, otp)
	if err != nil {
		return fmt.Errorf("could not validate otp: %w", err)
	}

	c.sessionToken = sessionToken
	c.refreshToken = refreshToken

	if err = c.sessionToken.WriteToFile(); err != nil {
		return fmt.Errorf("could not save token into file: %w", err)
	}

	if err = c.refreshToken.WriteToFile(); err != nil {
		return fmt.Errorf("could not save token into file: %w", err)
	}

	return nil
}

func (c *Client) refreshSession() error {
	c.logger.Info("refreshing session token")

	sessionToken, err := c.apiClient.Session(c.refreshToken)
	if err != nil {
		return fmt.Errorf("could not refresh session: %w", err)
	}

	c.sessionToken = sessionToken

	return nil
}

func (c *Client) SessionToken() api.Token {
	return c.sessionToken
}
