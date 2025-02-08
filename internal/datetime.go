package internal

import (
	"fmt"
	"time"
)

type DateTime struct {
	time.Time
}

// Formato de tiempo personalizado sin zona horaria
const customTimeFormat = "2006-01-02 15:04:05"

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format(customTimeFormat), nil
}

func (date *DateTime) UnmarshalCSV(csv string) error {
	t, err := time.Parse(customTimeFormat, csv)
	if err != nil {
		return fmt.Errorf("could not parse datetime: %w", err)
	}

	date.Time = t

	return nil
}
