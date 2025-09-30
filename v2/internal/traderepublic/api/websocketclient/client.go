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
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

const (
	// ConnectMsg is the message sent to establish a connection.
	ConnectMsg = "connect 31 {\"locale\":\"de\",\"platformId\":\"webtrading\"," +
		"\"platformVersion\":\"chrome - 134.0.0\",\"clientId\":\"app.traderepublic.com\",\"clientVersion\":\"3.174.0\"}"

	// Message types.
	MsgTypeSub   = "sub"
	MsgTypeUnsub = "unsub"

	// Message states.
	StateData     = traderepublic.WebsocketResponseSchemaJsonStateA
	StateContinue = traderepublic.WebsocketResponseSchemaJsonStateC
	StateError    = traderepublic.WebsocketResponseSchemaJsonStateE

	// Subscription types.
	TypeTimelineTransactions = "timelineTransactions"
	TypeTimelineDetail       = "timelineDetailV2"

	// Reconnect delay.
	reconnectDelay = 5 * time.Second

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
	publisher    publisher.Interface
	currentSubID uint
	mu           sync.Mutex
	closed       bool
	ctx          context.Context
}

// NewClient creates a new WebSocket client.
func NewClient(publisher publisher.Interface, ctx context.Context) *Client {
	client := &Client{
		publisher: publisher,
		ctx:       ctx,
	}

	client.Connect()

	return client
}

// Connect connects to the WebSocket server.
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	websocketURL := url.URL{Scheme: "wss", Host: internal.WebsocketBaseHost, Path: "/"}
	slog.Info("connecting to WebSocket", "url", websocketURL.String())

	// Create header with user agent
	header := make(map[string][]string)
	header["User-Agent"] = []string{internal.HTTPUserAgent}

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.DialContext(c.ctx, websocketURL.String(), header)
	if err != nil {
		return fmt.Errorf("could not connect to websocket: %w", err)
	}

	c.conn = conn
	c.closed = false

	// Send connect message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(ConnectMsg)); err != nil {
		return fmt.Errorf("could not send connect message: %w", err)
	}

	slog.Info("sent connect message")

	// Read the response
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("could not read connect response: %w", err)
	}

	slog.Info("received connect response", "response", string(msg))

	// Start goroutine to read messages
	go c.readMessages()

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
func (c *Client) Subscribe(data traderepublic.WebsocketSubRequestSchemaJson) (<-chan []byte, error) {
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

	ch := c.publisher.Subscribe(subID)

	// Create subscription message
	msg := fmt.Sprintf("%s %s %s", MsgTypeSub, subID, string(dataBytes))

	// Send subscription message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return nil, fmt.Errorf("could not send subscription message: %w", err)
	}

	slog.Debug("sent subscription message", "message", msg)

	return ch, nil
}

// readMessages reads messages from the WebSocket and sends them to the channel.
func (c *Client) readMessages() {
	for {
		select {
		case <-c.ctx.Done():
			slog.Info("context done, stopping message reader")

			err := c.Close()
			if err != nil {
				slog.Error("error closing connection", "error", err)
			}

			return
		default:
			// Read message
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				slog.Error("error reading message", "error", err)

				continue
			}

			// Parse message
			message, err := parseMessage(msg)
			if err != nil {
				slog.Error("error parsing message", "error", err, "message", string(msg))

				continue
			}

			// Handle message based on state
			switch message.State {
			case StateData:
				c.unsubscribe(message.ID)

				subID := strconv.FormatInt(int64(message.ID), 10)

				c.publisher.Publish([]byte(message.Data), subID)
				c.publisher.Close(subID)

				continue

			case StateContinue:
				// Continue reading
				slog.Debug("received continue message", "message", string(msg))

				continue

			case StateError:
				slog.Error("received error message", "message", string(msg))

				continue
			}
		}
	}
}

// unsubscribe unsubscribes from a subscription.
func (c *Client) unsubscribe(subID int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil || c.closed {
		return
	}

	// Create unsubscribe message
	msg := fmt.Sprintf("%s %d", MsgTypeUnsub, subID)

	// Send unsubscribe message
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		slog.Error("could not send unsubscribe message", "error", err)

		return
	}

	slog.Debug("sent unsubscribe message", "message", msg)
}

// reconnect reconnects to the WebSocket server.
func (c *Client) reconnect() error {
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
	slog.Info("reconnecting to WebSocket", "url", websocketURL.String())

	// Create header with user agent
	header := make(map[string][]string)
	header["User-Agent"] = []string{internal.HTTPUserAgent}

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.DialContext(c.ctx, websocketURL.String(), header)
	if err != nil {
		return fmt.Errorf("could not reconnect to websocket: %w", err)
	}

	c.conn = conn

	// Send connect message
	if err = c.conn.WriteMessage(websocket.TextMessage, []byte(ConnectMsg)); err != nil {
		return fmt.Errorf("could not send connect message: %w", err)
	}

	slog.Info("sent connect message")

	// Read the response
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("could not read connect response: %w", err)
	}

	slog.Info("received connect response", "response", string(msg))

	return nil
}

// parseMessage parses a message from the WebSocket.
func parseMessage(data []byte) (traderepublic.WebsocketResponseSchemaJson, error) {
	var msg traderepublic.WebsocketResponseSchemaJson

	parts := strings.Split(string(data), " ")

	if len(parts) < minMessageParts {
		return msg, errors.New("invalid message format")
	}

	// Parse ID
	id, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return msg, fmt.Errorf("could not parse message ID: %w", err)
	}

	msg.ID = int(id)
	msg.State = traderepublic.WebsocketResponseSchemaJsonState(parts[1])

	// Parse data if available
	if len(parts) > minMessageParts {
		msg.Data = strings.Join(parts[minMessageParts:], " ")
	}

	return msg, nil
}
