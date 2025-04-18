package websocketclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
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
	TypeTimeline   = "timeline"
	TypePortfolio  = "portfolio"
	TypeInstrument = "instrument"

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
func WithLogger(logger *log.Logger) ClientOption {
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
	logger       *log.Logger
	subID        uint
	mu           sync.Mutex
	closed       bool
}

// NewClient creates a new WebSocket client.
func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{
		logger: log.New(log.Writer(), "websocket-client: ", log.LstdFlags),
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
	c.logger.Printf("connecting to %s", websocketURL.String())

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

	c.logger.Println("sent connect message")

	// Read the response
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("could not read connect response: %w", err)
	}

	c.logger.Printf("received connect response: %s", string(msg))

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

// SubscribeToTimeline subscribes to timeline data.
func (c *Client) SubscribeToTimeline(ctx context.Context) (<-chan []byte, error) {
	data := map[string]interface{}{
		"type":  TypeTimeline,
		"token": c.sessionToken,
	}

	return c.subscribe(ctx, data)
}

// SubscribeToPortfolio subscribes to portfolio data.
func (c *Client) SubscribeToPortfolio(ctx context.Context) (<-chan []byte, error) {
	data := map[string]interface{}{
		"type":  TypePortfolio,
		"token": c.sessionToken,
	}

	return c.subscribe(ctx, data)
}

// SubscribeToInstrument subscribes to instrument data.
func (c *Client) SubscribeToInstrument(ctx context.Context, instrumentID string) (<-chan []byte, error) {
	data := map[string]interface{}{
		"type":  TypeInstrument,
		"id":    instrumentID,
		"token": c.sessionToken,
	}

	return c.subscribe(ctx, data)
}

// subscribe subscribes to a data type.
func (c *Client) subscribe(ctx context.Context, data map[string]interface{}) (<-chan []byte, error) {
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

	c.logger.Printf("sent subscription message: %s", msg)

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
			c.logger.Println("context done, stopping message reader")

			return
		default:
			// Read message
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				c.logger.Printf("error reading message: %v", err)

				// Try to reconnect
				if err := c.reconnect(ctx); err != nil {
					c.logger.Printf("could not reconnect: %v", err)

					return
				}

				continue
			}

			// Parse message
			message, err := parseMessage(msg)
			if err != nil {
				c.logger.Printf("error parsing message: %v", err)

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
					c.logger.Println("channel full, dropping message")
				}

				// Unsubscribe after receiving data
				c.unsubscribe(subID)

				return

			case StateContinue:
				// Continue reading
				c.logger.Println("received continue message")

			case StateError:
				// Parse error
				var errorResp ErrorResponse
				if err := json.Unmarshal(message.Data, &errorResp); err != nil {
					c.logger.Printf("error parsing error response: %v", err)

					continue
				}

				// Handle error
				for _, errorDetail := range errorResp.Errors {
					if errorDetail.IsUnauthorizedError() {
						c.logger.Println("unauthorized error, session expired")

						return
					}
				}

				c.logger.Printf("received error message: %s", string(message.Data))

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
	data := map[string]interface{}{
		"token": c.sessionToken,
	}

	// Marshal data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		c.logger.Printf("could not marshal unsubscribe data: %v", err)

		return
	}

	// Create unsubscribe message
	msg := fmt.Sprintf("%s %d %s", MsgTypeUnsub, subID, string(dataBytes))

	// Send unsubscribe message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		c.logger.Printf("could not send unsubscribe message: %v", err)

		return
	}

	c.logger.Printf("sent unsubscribe message: %s", msg)
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
	c.logger.Printf("reconnecting to %s", websocketURL.String())

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

	c.logger.Println("sent connect message")

	// Read the response
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("could not read connect response: %w", err)
	}

	c.logger.Printf("received connect response: %s", string(msg))

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
