package activitylog

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewClient,

	wire.Bind(new(ClientInterface), new(Client)),
)
