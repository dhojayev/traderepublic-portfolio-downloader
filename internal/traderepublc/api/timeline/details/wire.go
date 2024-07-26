package details

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewClient,
	NewTypeResolver,

	wire.Bind(new(ClientInterface), new(Client)),
	wire.Bind(new(TypeResolverInterface), new(TypeResolver)),
)

var TransactionSet = wire.NewSet(
	NewTransactionResponseNormalizer,

	wire.Bind(new(ResponseNormalizerInterface), new(TransactionResponseNormalizer)),
)

var ActivityLogSet = wire.NewSet(
	NewActivityResponseNormalizer,

	wire.Bind(new(ResponseNormalizerInterface), new(ActivityLogResponseNormalizer)),
)
