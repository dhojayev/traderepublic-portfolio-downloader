package filesystem

import (
	"github.com/google/wire"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

var CSVSet = wire.NewSet(
	NewCSVReader,
	NewCSVWriter,

	wire.Bind(new(CSVReaderInterface), new(CSVReader)),
	wire.Bind(new(CSVWriterInterface), new(CSVWriter)),
)

var JSONWriterSet = wire.NewSet(
	NewJSONWriter,

	wire.Bind(new(writer.Interface), new(*JSONWriter)),
)
