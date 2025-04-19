//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=client_interface.go -destination client_mock.go -package=websocketclient

package websocketclient

import "context"

// ClientInterface is the interface for the WebSocket client.
type ClientInterface interface {
	// Connect connects to the WebSocket server.
	Connect(ctx context.Context) error

	// Close closes the WebSocket connection.
	Close() error

	// SubscribeToTimelineTransactions subscribes to timeline transactions data.
	SubscribeToTimelineTransactions(ctx context.Context) (<-chan []byte, error)

	// SubscribeToTimelineTransactionsWithCursor subscribes to timeline transactions data with a cursor.
	SubscribeToTimelineTransactionsWithCursor(ctx context.Context, cursor string) (<-chan []byte, error)

	// SubscribeToTimelineDetail subscribes to timeline detail data.
	SubscribeToTimelineDetail(ctx context.Context, itemID string) (<-chan []byte, error)
}
