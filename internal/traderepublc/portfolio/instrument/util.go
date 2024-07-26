package instrument

import (
	"errors"
	"regexp"
)

var ErrNoMatch = errors.New("value did not match the pattern")

func ExtractInstrumentISINFromIcon(src string) (string, error) {
	pattern := regexp.MustCompile(`.*[^/]/([A-Z]{2}.*)/.*`)
	matches := pattern.FindStringSubmatch(src)

	if len(matches) == 0 {
		return "", ErrNoMatch
	}

	return matches[1], nil
}
