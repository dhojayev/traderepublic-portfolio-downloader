package traderepublic

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error constants for section and data item not found.
var (
	ErrSectionNotFound  = errors.New("section not found")
	ErrDataItemNotFound = errors.New("data item not found")

	// Title maps for payment details.
	PaymentShares     titleMap = titleMap{"Anteile"}
	PaymentSharePrice titleMap = titleMap{"Anteilspreis"}
	PaymentCommission titleMap = titleMap{"Gebühr"}
	PaymentTotal      titleMap = titleMap{"Gesamt"}
)

// titleMap is a slice of strings representing possible titles.
type titleMap []string

// SectionHeader retrieves the header section from the timeline details.
func (d *TimelineDetailsJson) SectionHeader() (HeaderSection, error) {
	var header HeaderSection

	// Find the header section in the sections slice.
	err := findSliceElement(d.Sections, &header, "")
	if err != nil {
		return header, fmt.Errorf("header %w", ErrSectionNotFound)
	}

	return header, nil
}

// SectionOverview retrieves the overview section from the timeline details.
func (d *TimelineDetailsJson) SectionOverview() (TableSection, error) {
	var overview TableSection

	// Find the overview section in the sections slice.
	err := findSliceElement(d.Sections, &overview, "Übersicht")
	if err != nil {
		return overview, fmt.Errorf("overview %w", ErrSectionNotFound)
	}

	return overview, nil
}

// SectionTransaction retrieves the transaction section from the timeline details.
func (d *TimelineDetailsJson) SectionTransaction() (TableSection, error) {
	var transaction TableSection

	// Find the transaction section in the sections slice.
	err := findSliceElement(d.Sections, &transaction, "Transaktion")
	if err != nil {
		return transaction, fmt.Errorf("transaction %w", ErrSectionNotFound)
	}

	return transaction, nil
}

// DataPayment retrieves a payment row based on the provided titles from the table section.
func (s *TableSection) DataPayment(titles titleMap) (PaymentRow, error) {
	var item PaymentRow

	// Iterate through the data slice to find the matching payment row.
	for _, title := range titles {
		err := findSliceElement(s.Data, &item, title)
		if err != nil {
			continue
		}

		return item, nil
	}

	return item, fmt.Errorf("payment %w with titles %#v", ErrDataItemNotFound, titles)
}

// findSliceElement searches for a slice element that matches the provided search criteria.
func findSliceElement(input []any, v any, search string) error {
	for _, element := range input {
		err := unmarshal(element, v)
		if err != nil {
			continue
		}

		// If no search criteria is provided, return the first match.
		if search == "" {
			return nil
		}

		// Check if the title matches the search criteria.
		title, ok := element.(map[string]any)["title"]
		if !ok {
			continue
		}

		if title != search {
			continue
		}

		return nil
	}

	return ErrSectionNotFound
}

// unmarshal converts an interface to a JSON string and then back to the provided value.
func unmarshal(i, v any) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}
