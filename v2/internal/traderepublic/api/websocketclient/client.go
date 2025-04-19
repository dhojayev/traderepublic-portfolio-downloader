package websocketclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
)

const (
	// ConnectMsg is the message sent to establish a connection.
	ConnectMsg = "connect 31 {\"locale\":\"de\",\"platformId\":\"webtrading\"," +
		"\"platformVersion\":\"chrome - 134.0.0\",\"clientId\":\"app.traderepublic.com\",\"clientVersion\":\"3.174.0\"}"

	// Message types.
	MsgTypeSub   = "sub"
	MsgTypeUnsub = "unsub"

	// Message states.
	StateData     = "A"
	StateContinue = "C"
	StateError    = "E"

	// Subscription types.
	TypeTimelineTransactions = "timelineTransactions"
	TypeTimelineDetail       = "timelineDetailV2"

	// Reconnect delay.
	reconnectDelay = 5 * time.Second

	// Channel buffer size.
	channelBufferSize = 10

	// Minimum parts in a message.
	minMessageParts = 2
)

var (
	// ErrNotConnected is returned when trying to use the client before connecting.
	ErrNotConnected = errors.New("websocket client not connected")

	// ErrConnectionClosed is returned when the connection is closed.
	ErrConnectionClosed = errors.New("websocket connection closed")

	// ErrAuthRequired is returned when authentication is required.
	ErrAuthRequired = errors.New("authentication required")
)

// Message represents a message received from the WebSocket.
type Message struct {
	ID    uint
	State string
	Data  []byte
}

// ErrorResponse represents an error response from the WebSocket.
type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}

// ErrorDetail represents a single error detail.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// IsAuthError returns true if the error is an authentication error.
func (e ErrorDetail) IsAuthError() bool {
	return e.Code == "AUTH_REQUIRED"
}

// IsUnauthorizedError returns true if the error is an unauthorized error.
func (e ErrorDetail) IsUnauthorizedError() bool {
	return e.Code == "UNAUTHORIZED"
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithLogger sets the logger for the client.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithSessionToken sets the session token for the client.
func WithSessionToken(token string) ClientOption {
	return func(c *Client) {
		c.sessionToken = token
	}
}

// Client is a WebSocket client for the Trade Republic API.
type Client struct {
	conn         *websocket.Conn
	sessionToken string
	logger       *slog.Logger
	subID        uint
	mu           sync.Mutex
	closed       bool
}

// NewClient creates a new WebSocket client.
func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

// Connect connects to the WebSocket server.
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	websocketURL := url.URL{Scheme: "wss", Host: internal.WebsocketBaseHost, Path: "/"}
	c.logger.Info("connecting to WebSocket", "url", websocketURL.String())

	// Create header with user agent
	header := make(map[string][]string)
	header["User-Agent"] = []string{internal.HTTPUserAgent}

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, websocketURL.String(), header)
	if err != nil {
		return fmt.Errorf("could not connect to websocket: %w", err)
	}

	c.conn = conn
	c.closed = false

	// Send connect message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(ConnectMsg)); err != nil {
		return fmt.Errorf("could not send connect message: %w", err)
	}

	c.logger.Info("sent connect message")

	// Read the response
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("could not read connect response: %w", err)
	}

	c.logger.Info("received connect response", "response", string(msg))

	return nil
}

// Close closes the WebSocket connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}

	c.closed = true

	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("could not close websocket connection: %w", err)
	}

	return nil
}

// prepareSubscription prepares a subscription request with the given parameters.
func (c *Client) prepareSubscription(dataType string, params map[string]any) map[string]any {
	data := map[string]any{
		"type":  dataType,
		"token": c.sessionToken,
	}

	// Add additional parameters
	for k, v := range params {
		data[k] = v
	}

	return data
}

// SubscribeToTimelineTransactions subscribes to timeline transactions data.
func (c *Client) SubscribeToTimelineTransactions(ctx context.Context) (<-chan []byte, error) {
	return c.SubscribeToTimelineTransactionsWithCursor(ctx, "")
}

// SubscribeToTimelineTransactionsWithCursor subscribes to timeline transactions data with a cursor.
func (c *Client) SubscribeToTimelineTransactionsWithCursor(ctx context.Context, cursor string) (<-chan []byte, error) {
	params := map[string]any{}

	// Add cursor if provided
	if cursor != "" {
		params["after"] = cursor
	}

	data := c.prepareSubscription(TypeTimelineTransactions, params)

	return c.subscribe(ctx, data)
}

// SubscribeToTimelineDetail subscribes to timeline detail data.
func (c *Client) SubscribeToTimelineDetail(ctx context.Context, itemID string) (<-chan []byte, error) {
	data := c.prepareSubscription(TypeTimelineDetail, map[string]any{"id": itemID})

	return c.subscribe(ctx, data)
}

// subscribe subscribes to a data type.
func (c *Client) subscribe(ctx context.Context, data map[string]any) (<-chan []byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil, ErrNotConnected
	}

	if c.closed {
		return nil, ErrConnectionClosed
	}

	if c.sessionToken == "" {
		return nil, ErrAuthRequired
	}

	c.subID++
	subID := c.subID

	// Add token to data
	data["token"] = c.sessionToken

	// Marshal data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data: %w", err)
	}

	// Create subscription message
	msg := fmt.Sprintf("%s %d %s", MsgTypeSub, subID, string(dataBytes))

	// Send subscription message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return nil, fmt.Errorf("could not send subscription message: %w", err)
	}

	c.logger.Info("sent subscription message", "message", msg)

	// Create channel for data
	dataCh := make(chan []byte, channelBufferSize)

	// Start goroutine to read messages
	go c.readMessages(ctx, subID, dataCh)

	return dataCh, nil
}

// readMessages reads messages from the WebSocket and sends them to the channel.
//
//nolint:gocognit,cyclop,funlen
func (c *Client) readMessages(ctx context.Context, subID uint, dataCh chan<- []byte) {
	defer close(dataCh)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("context done, stopping message reader")

			return
		default:
			// Read message
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				c.logger.Error("error reading message", "error", err)

				// Try to reconnect
				if err := c.reconnect(ctx); err != nil {
					c.logger.Error("could not reconnect", "error", err)

					return
				}

				continue
			}

			// Parse message
			message, err := parseMessage(msg)
			if err != nil {
				c.logger.Error("error parsing message", "error", err)

				continue
			}

			// Check if message is for this subscription
			if message.ID != subID {
				continue
			}

			// Handle message based on state
			switch message.State {
			case StateData:
				// Send data to channel
				select {
				case dataCh <- message.Data:
					// Data sent
				default:
					c.logger.Warn("channel full, dropping message")
				}

				// Unsubscribe after receiving data
				c.unsubscribe(subID)

				return

			case StateContinue:
				// Continue reading
				c.logger.Info("received continue message")

			case StateError:
				// Parse error
				var errorResp ErrorResponse
				if err := json.Unmarshal(message.Data, &errorResp); err != nil {
					c.logger.Error("error parsing error response", "error", err)

					continue
				}

				// Handle error
				for _, errorDetail := range errorResp.Errors {
					if errorDetail.IsUnauthorizedError() {
						c.logger.Warn("unauthorized error, session expired")

						return
					}
				}

				c.logger.Error("received error message", "data", string(message.Data))

				return
			}
		}
	}
}

// unsubscribe unsubscribes from a subscription.
func (c *Client) unsubscribe(subID uint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil || c.closed {
		return
	}

	// Create unsubscribe message
	data := map[string]any{
		"token": c.sessionToken,
	}

	// Marshal data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("could not marshal unsubscribe data", "error", err)

		return
	}

	// Create unsubscribe message
	msg := fmt.Sprintf("%s %d %s", MsgTypeUnsub, subID, string(dataBytes))

	// Send unsubscribe message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		c.logger.Error("could not send unsubscribe message", "error", err)

		return
	}

	c.logger.Info("sent unsubscribe message", "message", msg)
}

// reconnect reconnects to the WebSocket server.
func (c *Client) reconnect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}

	c.closed = false

	// Wait before reconnecting
	time.Sleep(reconnectDelay)

	// Connect to the WebSocket server
	websocketURL := url.URL{Scheme: "wss", Host: internal.WebsocketBaseHost, Path: "/"}
	c.logger.Info("reconnecting to WebSocket", "url", websocketURL.String())

	// Create header with user agent
	header := make(map[string][]string)
	header["User-Agent"] = []string{internal.HTTPUserAgent}

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, websocketURL.String(), header)
	if err != nil {
		return fmt.Errorf("could not reconnect to websocket: %w", err)
	}

	c.conn = conn

	// Send connect message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(ConnectMsg)); err != nil {
		return fmt.Errorf("could not send connect message: %w", err)
	}

	c.logger.Info("sent connect message")

	// Read the response
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("could not read connect response: %w", err)
	}

	c.logger.Info("received connect response", "response", string(msg))

	return nil
}

// parseMessage parses a message from the WebSocket.
func parseMessage(data []byte) (Message, error) {
	msg := Message{}
	parts := strings.Split(string(data), " ")

	if len(parts) < minMessageParts {
		return msg, errors.New("invalid message format")
	}

	// Parse ID
	id, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return msg, fmt.Errorf("could not parse message ID: %w", err)
	}

	msg.ID = uint(id)
	msg.State = parts[1]

	// Parse data if available
	if len(parts) > minMessageParts {
		msg.Data = []byte(strings.Join(parts[minMessageParts:], " "))
	}

	return msg, nil
}
