package traderepublic

import (
	"encoding/json"
	"errors"
)

var (
	ErrSectionNotFound  = errors.New("section not found")
	ErrDataItemNotFound = errors.New("data item not found")
)

func (d *TimelineDetailsJson) Section(v any) error {
	for _, section := range d.Sections {
		err := unmarshal(section, v)
		if err != nil {
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
