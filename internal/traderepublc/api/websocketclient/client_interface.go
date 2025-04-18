package websocketclient

import "context"

// ClientInterface is the interface for the WebSocket client.
type ClientInterface interface {
	// Connect connects to the WebSocket server
	Connect(ctx context.Context) error

	// Close closes the WebSocket connection
	Close() error

	// SubscribeToTimeline subscribes to timeline data
	SubscribeToTimeline(ctx context.Context) (<-chan []byte, error)

	// SubscribeToPortfolio subscribes to portfolio data
	SubscribeToPortfolio(ctx context.Context) (<-chan []byte, error)

	// SubscribeToInstrument subscribes to instrument data
	SubscribeToInstrument(ctx context.Context, instrumentID string) (<-chan []byte, error)
}
