package details

import (
	"errors"
	"fmt"
	"slices"
)

const (
	ResponseTimeFormat                      = "2006-01-02T15:04:05-0700"
	ResponseTimeFormatAlt                   = "2006-01-02T15:04:05.000000-07:00"
	ResponseSectionTypeValueHeader          = "header"
	ResponseSectionTypeValueTable           = "table"
	ResponseSectionTypeValueHorizontalTable = "horizontalTable"
	ResponseSectionTypeValueDocuments       = "documents"
	SectionTitleOverview                    = "Übersicht"
	SectionTitlePerformance                 = "Performance"
	SectionTitleTransaction                 = "Transaktion"
	SectionTitleTransactionAlt              = "Geschäft"
	SectionTitleSavingPlan                  = "Sparplan"
	OverviewDataTitleOrderType              = "Orderart"
	OverviewDataTitleAsset                  = "Asset"
	OverviewDataTitleUnderlyingAsset        = "Basiswert"
	OverviewDataTitleSecurity               = "Wertpapier"
	TransactionDataTitleShares              = "Anteile"
	TransactionDataTitleSharesAlt           = "Aktien"
	TransactionDataTitleRate                = "Aktienkurs"
	TransactionDataTitleRateAlt             = "Anteilspreis"
	TransactionDataTitleRateAlt2            = "Dividende je Aktie"
	TransactionDataTitleRateAlt3            = "Dividende pro Aktie"
	TransactionDataTitleCommission          = "Gebühr"
	TransactionDataTitleTotal               = "Gesamt"
	TransactionDataTitleTotalAlt            = "Summe"
	TransactionDataTitleTax                 = "Steuern"
	PerformanceDataTitleYield               = "Rendite"
	PerformanceDataTitleProfit              = "Gewinn"
	PerformanceDataTitleLoss                = "Verlust"
	OrderTypeTextsSale                      = "Verkauf"
	OrderTypeTextsPurchase                  = "Kauf"
	TrendNegative                           = "negative"
)

var ErrSectionDataTitleNotFound = errors.New("section data title not found")

type Response struct {
	ID       string           `json:"id"`
	Sections []map[string]any `json:"sections"`
}

type NormalizedResponse struct {
	ID          string
	Header      NormalizedResponseHeaderSection
	Overview    NormalizedResponseOverviewSection
	Performance NormalizedResponsePerformanceSection
	Transaction NormalizedResponseTransactionSection
	Documents   NormalizedResponseDocumentsSection
}

type NormalizedResponseHeaderSection struct {
	Action NormalizedResponseSectionAction     `json:"action"`
	Data   NormalizedResponseHeaderSectionData `json:"data"`
	Title  string                              `json:"title"`
	Type   string                              `json:"type"`
}

type NormalizedResponseTableSection struct {
	Data  []NormalizedResponseTableSectionData `json:"data"`
	Title string                               `json:"title"`
	Type  string                               `json:"type"`
}

func (s NormalizedResponseTableSection) GetDataByTitles(titles ...string) (NormalizedResponseTableSectionData, error) {
	for _, data := range s.Data {
		if !slices.Contains(titles, data.Title) {
			continue
		}

		return data, nil
	}

	return NormalizedResponseTableSectionData{}, fmt.Errorf("%w (%v)", ErrSectionDataTitleNotFound, titles)
}

type NormalizedResponseOverviewSection struct {
	NormalizedResponseTableSection
}

type NormalizedResponsePerformanceSection struct {
	NormalizedResponseTableSection
}

type NormalizedResponseTransactionSection struct {
	NormalizedResponseTableSection
}

type NormalizedResponseDocumentsSection struct {
	Data  []NormalizedResponseDocumentsSectionData `json:"data"`
	Title string                                   `json:"title"`
	Type  string                                   `json:"type"`
}

type NormalizedResponseHeaderSectionData struct {
	Icon         string `json:"icon"`
	Status       string `json:"status"`
	SubtitleText string `json:"subtitleText"`
	Timestamp    string `json:"timestamp"`
}

type NormalizedResponseTableSectionData struct {
	Detail NormalizedResponseTableSectionDataDetail `json:"detail"`
	Style  string                                   `json:"style"`
	Title  string                                   `json:"title"`
}

type NormalizedResponseTableSectionDataDetail struct {
	Action          NormalizedResponseSectionAction `json:"action"`
	FunctionalStyle string                          `json:"functionalStyle"`
	Amount          string                          `json:"amount"`
	Icon            string                          `json:"icon"`
	Status          string                          `json:"status"`
	Subtitle        string                          `json:"subtitle"`
	Timestamp       string                          `json:"timestamp"`
	Title           string                          `json:"title"`
	Text            string                          `json:"text"`
	Trend           string                          `json:"trend"`
	Type            string                          `json:"type"`
}

type NormalizedResponseDocumentsSectionData struct {
	Action      NormalizedResponseSectionAction `json:"action"`
	Detail      string                          `json:"detail"`
	ID          string                          `json:"id"`
	PostboxType string                          `json:"postboxType"`
	Title       string                          `json:"title"`
}

type NormalizedResponseSectionAction struct {
	Payload any    `json:"payload"`
	Type    string `json:"type"`
}
