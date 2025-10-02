package traderepublic

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrSectionNotFound = errors.New("section not found")

func (d *TimelineDetailsJson) Section(v any) error {
	for _, section := range d.Sections {
		data, err := json.Marshal(section)
		if err != nil {
			continue
		}

		err = json.Unmarshal(data, &v)
		if err != nil {
			continue
		}

		return nil
	}

	return fmt.Errorf("header %w", ErrSectionNotFound)
}
