package internal

import (
	"fmt"
	"time"
)

type DateTime struct {
	time.Time
}

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format(time.RFC822), nil
}

func (date *DateTime) UnmarshalCSV(csv string) error {
	t, err := time.Parse(time.RFC822, csv)
	if err != nil {
		return fmt.Errorf("could not parse datetime: %w", err)
	}

	date.Time = t

	return nil
}
