//go:generate go tool mockgen -source=client_interface.go -destination client_mock.go -package=websocketclient

package websocketclient

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

// ClientInterface is the interface for the WebSocket client.
type ClientInterface interface {
	// Connect connects to the WebSocket server.
	Connect() error

	// Close closes the WebSocket connection.
	Close() error

	// Subscribe subscribes to a data type.
	Subscribe(data traderepublic.WsSubRequestJson) (<-chan []byte, error)
}
