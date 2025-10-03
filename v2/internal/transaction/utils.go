package transaction

import (
	"fmt"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
)

type CSVDateTime struct {
	time.Time
}

func (date *CSVDateTime) MarshalCSV() (string, error) {
	return date.Time.In(time.Local).Format(time.RFC822Z), nil
}

func (date *CSVDateTime) UnmarshalCSV(csv string) error {
	t, err := time.Parse(time.RFC822Z, csv)
	if err != nil {
		return fmt.Errorf("could not parse datetime: %w", err)
	}

	date.Time = t

	return nil
}

func ParseTimestamp(src string) (time.Time, error) {
	timestamp, err := time.Parse(internal.ResponseTimeFormat, src)

	if err == nil {
		return timestamp, nil
	}

	timestamp, err2 := time.Parse(internal.ResponseTimeFormatAlt, src)

	if err2 == nil {
		return timestamp, nil
	}

	return time.Time{}, fmt.Errorf("could not parse timestamp 2 times: '%w' and '%w'", err, err2)
}
