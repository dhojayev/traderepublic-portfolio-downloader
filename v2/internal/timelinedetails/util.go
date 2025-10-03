package timelinedetails

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ErrNoMatch = errors.New("value did not match the pattern")

func ParseFloatFromResponse(src string) (float64, error) {
	pattern := regexp.MustCompile(`(\d+(?:\.\d+)*|\d+)(?:,(\d+))?`)
	matches := pattern.FindStringSubmatch(src)

	if len(matches) == 0 {
		return 0, ErrNoMatch
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
