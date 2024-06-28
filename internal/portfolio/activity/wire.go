package activity

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewProcessor,
	NewHandler,

	wire.Bind(new(ProcessorInterface), new(Processor)),
	wire.Bind(new(HandlerInterface), new(Handler)),
)
