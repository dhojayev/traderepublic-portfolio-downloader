package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

var (
	TransactionTestCasesSupported   []TransactionTestCase
	TransactionTestCasesUnsupported []TransactionTestCase
	TransactionTestCasesUnknown     []TransactionTestCase
	ActivityLogTestCasesSupported   []ActivityLogTestCase
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
	Raw          []byte
	Unmarshalled TimelineDetailsResponseSections
}

type TimelineDetailsResponseSections struct {
	Header    details.ResponseSectionTypeHeader
	Table     details.ResponseSectionsTypeTable
	Documents details.ResponseSectionTypeDocuments
}

func RegisterSupported(testCase TransactionTestCase) {
	TransactionTestCasesSupported = append(TransactionTestCasesSupported, testCase)
}

func RegisterUnsupported(testCase TransactionTestCase) {
	TransactionTestCasesUnsupported = append(TransactionTestCasesUnsupported, testCase)
}

func RegisterUnknown(testCase TransactionTestCase) {
	TransactionTestCasesUnknown = append(TransactionTestCasesUnknown, testCase)
}

func RegisterActivityLogSupported(testCase ActivityLogTestCase) {
	ActivityLogTestCasesSupported = append(ActivityLogTestCasesSupported, testCase)
}
