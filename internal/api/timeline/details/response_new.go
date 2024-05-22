package details

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TODO add support of: TAX_REFUND

const (
	ResponseSectionTypeValueHeader    = "header"
	ResponseSectionTypeValueTable     = "table"
	ResponseSectionTypeValueDocuments = "documents"
)

var (
	ErrSectionContainsNoType = errors.New("section contains no type")
	ErrSectionTypeNotFound   = errors.New("section type not found")
)

type ResponseNew struct {
	ID       string           `json:"id"`
	Sections []map[string]any `json:"sections"`
}

func (r ResponseNew) SectionTypeHeader() (ResponseSectionTypeHeaderNew, error) {
	var sections []ResponseSectionTypeHeaderNew

	if err := r.deserializeSections(ResponseSectionTypeValueHeader, &sections); err != nil {
		return ResponseSectionTypeHeaderNew{}, err
	}

	if len(sections) == 0 {
		return ResponseSectionTypeHeaderNew{}, ErrSectionTypeNotFound
	}

	return sections[0], nil
}

func (r ResponseNew) SectionsTypeTable() ([]ResponseSectionTypeTableNew, error) {
	var sections []ResponseSectionTypeTableNew

	if err := r.deserializeSections(ResponseSectionTypeValueTable, &sections); err != nil {
		return sections, err
	}

	return sections, nil
}

func (r ResponseNew) SectionTypeDocuments() (ResponseSectionTypeDocumentsNew, error) {
	var sections []ResponseSectionTypeDocumentsNew

	if err := r.deserializeSections(ResponseSectionTypeValueDocuments, &sections); err != nil {
		return ResponseSectionTypeDocumentsNew{}, err
	}

	if len(sections) == 0 {
		return ResponseSectionTypeDocumentsNew{}, ErrSectionTypeNotFound
	}

	return sections[0], nil
}

func (r ResponseNew) deserializeSections(needle string, v any) error {
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

type ResponseSectionTypeHeaderNew struct {
	Action ResponseActionNew                `json:"action"`
	Data   ResponseSectionTypeHeaderDataNew `json:"data"`
	Title  string                           `json:"title"`
	Type   string                           `json:"type"`
}

type ResponseSectionTypeTableNew struct {
	// Action ResponseActionNew                 `json:"action"`
	Data  []ResponseSectionTypeTableDataNew `json:"data"`
	Title string                            `json:"title"`
	Type  string                            `json:"type"`
}

type ResponseSectionTypeDocumentsNew struct {
	Action ResponseActionNew                    `json:"action"`
	Data   []ResponseSectionTypeDocumentDataNew `json:"data"`
	Title  string                               `json:"title"`
	Type   string                               `json:"type"`
}

type ResponseActionNew struct {
	Payload any    `json:"payload"`
	Type    string `json:"type"`
}

type ResponseSectionTypeHeaderDataNew struct {
	Icon         string `json:"icon"`
	Status       string `json:"status"`
	SubtitleText string `json:"subtitleText"`
	Timestamp    string `json:"timestamp"`
}

type ResponseSectionTypeTableDataNew struct {
	Detail ResponseSectionTypeTableDataDetailNew `json:"detail"`
	Style  string                                `json:"style"`
	Title  string                                `json:"title"`
}

type ResponseSectionTypeDocumentDataNew struct {
	Action      ResponseActionNew `json:"action"`
	Detail      string            `json:"detail"`
	ID          string            `json:"id"`
	PostboxType string            `json:"postboxType"`
	Title       string            `json:"title"`
}

type ResponseSectionTypeTableDataDetailNew struct {
	Action          ResponseActionNew `json:"action"`
	FunctionalStyle string            `json:"functionalStyle"`
	Amount          string            `json:"amount"`
	Icon            string            `json:"icon"`
	Status          string            `json:"status"`
	Subtitle        string            `json:"subtitle"`
	Timestamp       string            `json:"timestamp"`
	Title           string            `json:"title"`
	Text            string            `json:"text"`
	Trend           string            `json:"trend"`
	Type            string            `json:"type"`
}
