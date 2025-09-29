package websocketclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message/publisher"
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

// Client is a WebSocket client for the Trade Republic API.
type Client struct {
	conn         *websocket.Conn
	publisher    publisher.PublisherInterface
	logger       *slog.Logger
	currentSubID uint
	mu           sync.Mutex
	closed       bool
}

// NewClient creates a new WebSocket client.
func NewClient(publisher publisher.PublisherInterface, logger *slog.Logger) *Client {
	return &Client{
		publisher: publisher,
		logger:    logger,
	}
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

// subscribe subscribes to a data type.
func (c *Client) Subscribe(ctx context.Context, data map[string]any) (<-chan []byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil, ErrNotConnected
	}

	if c.closed {
		return nil, ErrConnectionClosed
	}

	c.currentSubID++
	subID := strconv.FormatUint(uint64(c.currentSubID), 10)

	// Marshal data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data: %w", err)
	}

	// Create subscription message
	msg := fmt.Sprintf("%s %s %s", MsgTypeSub, subID, string(dataBytes))

	// Send subscription message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return nil, fmt.Errorf("could not send subscription message: %w", err)
	}

	c.logger.Info("sent subscription message", "message", msg)

	sub := c.publisher.Subscribe(subID)

	// Start goroutine to read messages
	go c.readMessages(ctx, subID)

	return sub, nil
}

// readMessages reads messages from the WebSocket and sends them to the channel.
func (c *Client) readMessages(ctx context.Context, subID string) {
	defer c.publisher.Close(subID)

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

				return
			}

			// Parse message
			message, err := parseMessage(msg)
			if err != nil {
				c.logger.Error("error parsing message", "error", err)

				continue
			}

			// Handle message based on state
			switch message.State {
			case StateData:
				c.publisher.Publish(message.Data, string(subID))
				c.unsubscribe(subID, "")

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
func (c *Client) unsubscribe(subID string, token string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil || c.closed {
		return
	}

	// Create unsubscribe message
	data := map[string]any{
		"token": token,
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
