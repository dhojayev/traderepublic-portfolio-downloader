package traderepublic

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrSectionNotFound  = errors.New("section not found")
	ErrDataItemNotFound = errors.New("data item not found")
)

func (d *TimelineDetailsJson) SectionHeader() (HeaderSection, error) {
	var header HeaderSection

	err := findSliceElement(d.Sections, &header, "")
	if err != nil {
		return header, fmt.Errorf("header %w", ErrSectionNotFound)
	}

	return header, nil
}

func (d *TimelineDetailsJson) SectionOverview() (TableSection, error) {
	var overview TableSection

	err := findSliceElement(d.Sections, &overview, "Übersicht")
	if err != nil {
		return overview, fmt.Errorf("overview %w", ErrSectionNotFound)
	}

	return overview, nil
}

func (d *TimelineDetailsJson) SectionTransaction() (TableSection, error) {
	var transaction TableSection

	err := findSliceElement(d.Sections, &transaction, "Transaktion")
	if err != nil {
		return transaction, fmt.Errorf("transaction %w", ErrSectionNotFound)
	}

	return transaction, nil
}

func (s *TableSection) DataShares() (PaymentRow, error) {
	var shares PaymentRow

	err := findSliceElement(s.Data, &shares, "Anteile")
	if err != nil {
		return shares, fmt.Errorf("shares %w", ErrDataItemNotFound)
	}

	return shares, nil
}

func findSliceElement(input []any, v any, search string) error {
	for _, element := range input {
		err := unmarshal(element, v)
		if err != nil {
			continue
		}

		if search == "" {
			return nil
		}

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
