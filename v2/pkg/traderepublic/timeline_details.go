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

	SectionTableOverview    = sectionTableTitles{"Overview"}    // Title for the overview table section
	SectionTableTransaction = sectionTableTitles{"Transaction"} // Title for the transaction table section

	// Title maps for payment details.
	DataShares           = dataTitles{"Shares"}      // Title map for shares in payment details
	DataSharePrice       = dataTitles{"Share price"} // Title map for share price in payment details
	DataFee              = dataTitles{"Fee"}         // Title map for commission in payment details
	DataTotal            = dataTitles{"Total"}       // Title map for total in payment details
	DataTax              = dataTitles{"Tax"}
	DataDividendPerShare = dataTitles{"Dividend per share"}
)

// sectionTableTitles is a type alias for string representing a table section title.
type sectionTableTitles []string

// dataTitles is a slice of strings representing possible titles.
type dataTitles []string

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

func (d *TimelineDetailsJson) SectionTable(titles sectionTableTitles) (TableSection, error) {
	var section TableSection

	for _, title := range titles {
		// Find the section in the sections slice by title.
		err := findSliceElement(d.Sections, &section, string(title))
		if err != nil {
			continue
		}

		return section, nil
	}

	return section, fmt.Errorf("table %w with titles %#v", ErrSectionNotFound, titles)
}

// DataPayment retrieves a payment row based on the provided titles from the table section.
func (s *TableSection) DataPayment(titles dataTitles) (PaymentRow, error) {
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
