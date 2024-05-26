package tests

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

type TestCase struct {
	ResponseJSON string
	EventType    transactions.EventType
	Transaction  transaction.Model
	CSVEntry     filesystem.CSVEntry
}
