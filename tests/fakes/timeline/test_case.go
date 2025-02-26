package timeline_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
)

var TestCases = make(map[string]Fake)

type Fake struct {
	RawResponse  []byte
	Unmarshalled transactions.ResponseItem
}
