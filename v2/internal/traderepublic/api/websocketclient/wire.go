package websocketclient

import (
	"log/slog"

	"github.com/google/wire"
)

// ProvideClient provides a WebSocket client.
func ProvideClient(logger *slog.Logger, sessionToken string) (*Client, error) {
	return NewClient(
		WithLogger(logger),
		WithSessionToken(sessionToken),
	)
}

// Set is a wire.ProviderSet that provides the WebSocket client.
var Set = wire.NewSet(
	ProvideClient,
	wire.Bind(new(ClientInterface), new(*Client)),
)
