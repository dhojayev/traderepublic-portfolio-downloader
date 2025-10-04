package traderepublic

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error constants for section and data item not found.
var (
	ErrSliceElementNotFound = errors.New("slice element not found")
	ErrSectionNotFound      = errors.New("section not found")
	ErrDataItemNotFound     = errors.New("data item not found")
	ErrStepNotFound         = errors.New("step not found")

	SectionOverview    = sectionTitles{"Overview"}    // Title for the overview table section
	SectionTransaction = sectionTitles{"Transaction"} // Title for the transaction table section
	SectionPerformance = sectionTitles{"Performance"}
	SectionSender      = sectionTitles{"Sender"}

	StepInterestPayment = "Interest payment"

	DataFrom             = dataTitles{"From"}
	DataTo               = dataTitles{"To"}
	DataPayment          = dataTitles{"Payment"}
	DataRoundUp          = dataTitles{"Round up"}
	DataBuy              = dataTitles{"Buy"}
	DataCardVerification = dataTitles{"Card verification"}
	DataCardPayment      = dataTitles{"Card payment"}
	DataCardRefund       = dataTitles{"Card refund"}
	DataAverageBalance   = dataTitles{"Average balance"}
	DataSell             = dataTitles{"Sell"}
	DataLimitSell        = dataTitles{"Limit Sell"}
	DataSaveback         = dataTitles{"Saveback"}
	DataOrderType        = dataTitles{"Order Type"}
	DataSavingsPlan      = dataTitles{"Savings Plan"}
	DataEvent            = dataTitles{"Event"}
	DataAsset            = dataTitles{"Asset"}
	DataShares           = dataTitles{"Shares"}      // Title map for shares in payment details
	DataSharePrice       = dataTitles{"Share price"} // Title map for share price in payment details
	DataFee              = dataTitles{"Fee"}         // Title map for commission in payment details
	DataProfit           = dataTitles{"Profit"}
	DataGain             = dataTitles{"Gain"}
	DataTotal            = dataTitles{"Total"} // Title map for total in payment details
	DataTax              = dataTitles{"Tax"}
	DataDividendPerShare = dataTitles{"Dividend per share"}
)

// sectionTitles is a type alias for string representing a table section title.
type sectionTitles []string

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

func (d *TimelineDetailsJson) SectionSteps() (StepsSection, error) {
	var steps StepsSection

	err := findSliceElement(d.Sections, &steps, "")
	if err != nil {
		return steps, fmt.Errorf("steps %w", ErrSectionNotFound)
	}

	return steps, nil
}

func (d *TimelineDetailsJson) FindSection(titles sectionTitles) (TableSection, error) {
	var section TableSection

	for _, title := range titles {
		// Find the section in the sections slice by title.
		err := findSliceElement(d.Sections, &section, string(title))
		if err != nil {
			continue
		}

		return section, nil
	}

	return TableSection{}, fmt.Errorf("table %w with titles %v", ErrSectionNotFound, titles)
}

// FindData retrieves a payment row based on the provided titles from the table section.
func (s *TableSection) FindData(titles dataTitles) (PaymentRow, error) {
	var item PaymentRow

	// Iterate through the data slice to find the matching payment row.
	for _, title := range titles {
		err := findSliceElement(s.Data, &item, title)
		if err != nil {
			continue
		}

		return item, nil
	}

	return PaymentRow{}, fmt.Errorf("%w with titles %v", ErrDataItemNotFound, titles)
}

func (s *StepsSection) FindStep(title string) (StepItem, error) {
	for _, step := range s.Steps {
		if step.Content.Title != title {
			continue
		}

		return step, nil
	}

	return StepItem{}, fmt.Errorf("%w with title %s", ErrStepNotFound, title)
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

	return ErrSliceElementNotFound
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
