package writer

import "github.com/google/wire"

var NilSet = wire.NewSet(
	NewNilWriter,

	wire.Bind(new(Interface), new(NilWriter)),
)
