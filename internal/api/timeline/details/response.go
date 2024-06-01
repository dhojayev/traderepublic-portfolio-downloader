package details

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ResponseTimeFormat                      = "2006-01-02T15:04:05-0700"
	ResponseSectionTypeValueHeader          = "header"
	ResponseSectionTypeValueTable           = "table"
	ResponseSectionTypeValueHorizontalTable = "horizontalTable"
	ResponseSectionTypeValueDocuments       = "documents"
	SectionTitleOverview                    = "Übersicht"
	SectionTitlePerformance                 = "Performance"
	SectionTitleTransaction                 = "Transaktion"
	OverviewDataTitleOrderType              = "Orderart"
	OverviewDataTitleAsset                  = "Asset"
	OverviewDataTitleUnderlyingAsset        = "Basiswert"
	TransactionDataTitleShares              = "Anteile"
	TransactionDataTitleSharesAlt           = "Aktien"
	TransactionDataTitleRate                = "Aktienkurs"
	TransactionDataTitleRateAlt             = "Anteilspreis"
	TransactionDataTitleRateAlt2            = "Dividende je Aktie"
	TransactionDataTitleCommission          = "Gebühr"
	TransactionDataTitleTotal               = "Gesamt"
	TransactionDataTitleTax                 = "Steuern"
	PerformanceDataTitleYield               = "Rendite"
	PerformanceDataTitleProfit              = "Gewinn"
	PerformanceDataTitleLoss                = "Verlust"
	OrderTypeTextsSale                      = "Verkauf"
	OrderTypeTextsPurchase                  = "Kauf"
	TrendNegative                           = "negative"
)

var (
	ErrSectionContainsNoType    = errors.New("section contains no type")
	ErrSectionTypeNotFound      = errors.New("section type not found")
	ErrSectionTitleNotFound     = errors.New("section title not found")
	ErrSectionDataTitleNotFound = errors.New("section data title not found")
)

type Response struct {
	ID       string           `json:"id"`
	Sections []map[string]any `json:"sections"`
}

func (r Response) SectionTypeHeader() (ResponseSectionTypeHeader, error) {
	var sections []ResponseSectionTypeHeader

	if err := r.deserializeSections(ResponseSectionTypeValueHeader, &sections); err != nil {
		return ResponseSectionTypeHeader{}, err
	}

	if len(sections) == 0 {
		return ResponseSectionTypeHeader{}, ErrSectionTypeNotFound
	}

	return sections[0], nil
}

func (r Response) SectionsTypeTable() (ResponseSectionsTypeTable, error) {
	var sections ResponseSectionsTypeTable

	if err := r.deserializeSections(ResponseSectionTypeValueTable, &sections); err != nil {
		return sections, err
	}

	return sections, nil
}

func (r Response) SectionsTypeHorizontalTable() (ResponseSectionsTypeTable, error) {
	var sections ResponseSectionsTypeTable

	if err := r.deserializeSections(ResponseSectionTypeValueHorizontalTable, &sections); err != nil {
		return sections, err
	}

	return sections, nil
}

func (r Response) SectionTypeDocuments() (ResponseSectionTypeDocuments, error) {
	var sections []ResponseSectionTypeDocuments

	if err := r.deserializeSections(ResponseSectionTypeValueDocuments, &sections); err != nil {
		return ResponseSectionTypeDocuments{}, err
	}

	if len(sections) == 0 {
		return ResponseSectionTypeDocuments{}, ErrSectionTypeNotFound
	}

	return sections[0], nil
}

func (r Response) deserializeSections(needle string, v any) error {
	sections := make([]map[string]any, 0)

	for _, section := range r.Sections {
		sectionType, ok := section["type"]
		if !ok {
			return ErrSectionContainsNoType
		}

		if sectionType != needle {
			continue
		}

		sections = append(sections, section)
	}

	if len(sections) == 0 {
		return ErrSectionTypeNotFound
	}

	sectionsBytes, err := json.Marshal(sections)
	if err != nil {
		return fmt.Errorf("could not marshal %s section: %w", needle, err)
	}

	if err = json.Unmarshal(sectionsBytes, v); err != nil {
		return fmt.Errorf("could not unmarshal %s section: %w", needle, err)
	}

	return nil
}

type ResponseSectionTypeHeader struct {
	Action ResponseAction                `json:"action"`
	Data   ResponseSectionTypeHeaderData `json:"data"`
	Title  string                        `json:"title"`
	Type   string                        `json:"type"`
}

type ResponseSectionTypeTable struct {
	// Action ResponseAction                 `json:"action"`
	Data  []ResponseSectionTypeTableData `json:"data"`
	Title string                         `json:"title"`
	Type  string                         `json:"type"`
}

func (s ResponseSectionTypeTable) GetDataByTitle(title string) (ResponseSectionTypeTableData, error) {
	for _, data := range s.Data {
		if data.Title != title {
			continue
		}

		return data, nil
	}

	return ResponseSectionTypeTableData{}, fmt.Errorf("%w (%s)", ErrSectionDataTitleNotFound, title)
}

type ResponseSectionsTypeTable []ResponseSectionTypeTable

func (s ResponseSectionsTypeTable) FindByTitle(title string) (ResponseSectionTypeTable, error) {
	for _, section := range s {
		if section.Title != title {
			continue
		}

		return section, nil
	}

	return ResponseSectionTypeTable{}, fmt.Errorf("%w (%s)", ErrSectionTitleNotFound, title)
}

type ResponseSectionTypeDocuments struct {
	Action ResponseAction                    `json:"action"`
	Data   []ResponseSectionTypeDocumentData `json:"data"`
	Title  string                            `json:"title"`
	Type   string                            `json:"type"`
}

type ResponseAction struct {
	Payload any    `json:"payload"`
	Type    string `json:"type"`
}

type ResponseSectionTypeHeaderData struct {
	Icon         string `json:"icon"`
	Status       string `json:"status"`
	SubtitleText string `json:"subtitleText"`
	Timestamp    string `json:"timestamp"`
}

type ResponseSectionTypeTableData struct {
	Detail ResponseSectionTypeTableDataDetail `json:"detail"`
	Style  string                             `json:"style"`
	Title  string                             `json:"title"`
}

type ResponseSectionTypeDocumentData struct {
	Action      ResponseAction `json:"action"`
	Detail      string         `json:"detail"`
	ID          string         `json:"id"`
	PostboxType string         `json:"postboxType"`
	Title       string         `json:"title"`
}

type ResponseSectionTypeTableDataDetail struct {
	Action          ResponseAction `json:"action"`
	FunctionalStyle string         `json:"functionalStyle"`
	Amount          string         `json:"amount"`
	Icon            string         `json:"icon"`
	Status          string         `json:"status"`
	Subtitle        string         `json:"subtitle"`
	Timestamp       string         `json:"timestamp"`
	Title           string         `json:"title"`
	Text            string         `json:"text"`
	Trend           string         `json:"trend"`
	Type            string         `json:"type"`
}
