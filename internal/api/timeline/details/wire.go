package details

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewClient,
	NewTypeResolver,
	NewResponseNormalizer,

	wire.Bind(new(ClientInterface), new(Client)),
	wire.Bind(new(TypeResolverInterface), new(TypeResolver)),
	wire.Bind(new(ResponseNormalizerInterface), new(ResponseNormalizer)),
)
