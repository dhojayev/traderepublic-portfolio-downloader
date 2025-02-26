package details_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
	timeline_test "github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes/timeline"
)

var (
	TestCasesSupported   = make(map[string]Fake)
	TestCasesUnsupported = make(map[string]Fake)
	TestCasesUnknown     = make(map[string]Fake)
)

type Fake struct {
	Parent       *timeline_test.Fake
	RawResponse  []byte
	Unmarshalled details.NormalizedResponse
	Model        transaction.Model
	CSVEntry     filesystem.CSVEntry
}
