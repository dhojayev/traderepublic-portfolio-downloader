package transaction

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
)

var ErrPatternMismatch = errors.New("value did not match the pattern")

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

func ExtractInstrumentISINFromIcon(src string) (string, error) {
	pattern := regexp.MustCompile(`.*[^/]/([A-Z]{2}.*)/.*`)
	matches := pattern.FindStringSubmatch(src)

	if len(matches) == 0 {
		return "", ErrPatternMismatch
	}

	return matches[1], nil
}

func ParseFloatFromResponse(src string) (float64, error) {
	pattern := regexp.MustCompile(`(\d+(?:\.\d+)*|\d+)(?:,(\d+))?`)
	matches := pattern.FindStringSubmatch(src)

	if len(matches) == 0 {
		return 0, ErrPatternMismatch
	}

	wholePart := matches[1]
	decimalPart := matches[2]
	strFloat := wholePart

	if decimalPart != "" {
		strFloat = strings.ReplaceAll(wholePart, ".", "") + "." + decimalPart
	}

	value, err := strconv.ParseFloat(strFloat, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse float from '%s': %w", src, err)
	}

	return value, nil
}
