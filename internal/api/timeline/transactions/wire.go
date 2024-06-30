package transactions

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewClient,
	NewEventTypeResolver,

	wire.Bind(new(ClientInterface), new(Client)),
	wire.Bind(new(EventTypeResolverInterface), new(EventTypeResolver)),
)
