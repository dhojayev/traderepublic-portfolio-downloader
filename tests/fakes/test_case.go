package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

var (
	TransactionTestCasesSupported   = make(map[string]TransactionTestCase)
	TransactionTestCasesUnsupported = make(map[string]TransactionTestCase)
	TransactionTestCasesUnknown     = make(map[string]TransactionTestCase)
	ActivityLogTestCasesSupported   = make(map[string]ActivityLogTestCase)
)

type TransactionTestCase struct {
	TimelineTransactionsData TimelineTransactionsTestData
	TimelineDetailsData      TimelineDetailsTestData
	EventType                transactions.EventType
	Transaction              transaction.Model
	CSVEntry                 filesystem.CSVEntry
}

type ActivityLogTestCase struct {
	ActivityLogData     ActivityLogTestData
	TimelineDetailsData TimelineDetailsTestData
}

type TimelineTransactionsTestData struct {
	Raw          []byte
	Unmarshalled transactions.ResponseItem
}

type ActivityLogTestData struct {
	Raw          []byte
	Unmarshalled activitylog.ResponseItem
}

type TimelineDetailsTestData struct {
	Raw        []byte
	Normalized details.NormalizedResponse
}

func RegisterSupported(name string, testCase TransactionTestCase) {
	TransactionTestCasesSupported[name] = testCase
}

func RegisterUnsupported(name string, testCase TransactionTestCase) {
	TransactionTestCasesUnsupported[name] = testCase
}

func RegisterUnknown(name string, testCase TransactionTestCase) {
	TransactionTestCasesUnknown[name] = testCase
}

func RegisterActivityLogSupported(name string, testCase ActivityLogTestCase) {
	ActivityLogTestCasesSupported[name] = testCase
}
