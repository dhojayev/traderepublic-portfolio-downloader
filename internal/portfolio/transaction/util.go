package transaction

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var ErrNoMatch = errors.New("value did not match the pattern")

func ExtractInstrumentNameFromIcon(src string) (string, error) {
	pattern := regexp.MustCompile(`.*[^/]/(.*)/.*`)
	matches := pattern.FindStringSubmatch(src)

	if len(matches) == 0 {
		return "", ErrNoMatch
	}

	return matches[1], nil
}

func ParseFloatWithPeriod(src string) (float64, error) {
	pattern := regexp.MustCompile(`^(\d+)\.?(\d*)$`)
	matches := pattern.FindStringSubmatch(src)

	if len(matches) == 0 {
		return 0, ErrNoMatch
	}

	value := matches[1] + "." + matches[2]

	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse float from '%s': %w", value, err)
	}

	return valueFloat, nil
}

func ParseFloatWithComma(src string, isNegative bool) (float64, error) {
	pattern := regexp.MustCompile(`^\+?\s?(\d+)\.?(\d*),?(\d*).*$`)
	matches := pattern.FindStringSubmatch(src)

	if len(matches) == 0 {
		return 0, ErrNoMatch
	}

	value := matches[1] + matches[2] + "." + matches[3]

	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse float from '%s': %w", value, err)
	}

	if isNegative {
		valueFloat = -valueFloat
	}

	return valueFloat, nil
}
