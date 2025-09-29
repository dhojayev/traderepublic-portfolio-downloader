//go:generate go tool mockgen -source=client_interface.go -destination client_mock.go -package=websocketclient

package websocketclient

import "context"

// ClientInterface is the interface for the WebSocket client.
type ClientInterface interface {
	// Connect connects to the WebSocket server.
	Connect(ctx context.Context) error

	// Close closes the WebSocket connection.
	Close() error

	// Subscribe subscribes to a data type.
	Subscribe(ctx context.Context, data map[string]any) (<-chan []byte, error)
}
