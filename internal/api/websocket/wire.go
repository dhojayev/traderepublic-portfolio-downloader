package websocket

import (
	"github.com/google/wire"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
)

var DefaultSet = wire.NewSet(
	NewReader,

	wire.Bind(new(reader.Interface), new(*Reader)),
)
