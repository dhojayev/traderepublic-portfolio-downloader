package tests

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

var (
	TestCasesSupported   []TestCase
	TestCasesUnsupported []TestCase
	TestCasesUnknown     []TestCase
)

type TestCase struct {
	ResponseJSON string
	Response     Response
	EventType    transactions.EventType
	Transaction  transaction.Model
	CSVEntry     filesystem.CSVEntry
}

type Response struct {
	HeaderSection    details.ResponseSectionTypeHeader
	TableSections    details.ResponseSectionsTypeTable
	DocumentsSection details.ResponseSectionTypeDocuments
}

func RegisterSupported(testCase TestCase) {
	TestCasesSupported = append(TestCasesSupported, testCase)
}

func RegisterUnsupported(testCase TestCase) {
	TestCasesUnsupported = append(TestCasesUnsupported, testCase)
}

func RegisterUnknown(testCase TestCase) {
	TestCasesUnknown = append(TestCasesUnknown, testCase)
}
