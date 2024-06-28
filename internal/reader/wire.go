package reader

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewJSONReader,

	wire.Bind(new(Interface), new(*JSONReader)),
)
