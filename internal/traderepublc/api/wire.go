package api

import "github.com/google/wire"

// ClientSet is a wire set that uses the generated OpenAPI client.
var ClientSet = wire.NewSet(
	NewClient,
	NewWSClient,

	wire.Bind(new(ClientInterface), new(*Client)),
)
