package internal

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thlib/go-timezone-local/tzlocal"

	// load timezone data.
	_ "time/tzdata"
)

func GetRuntimeTimezone(logger *log.Logger) error {
	timezoneName, err := tzlocal.RuntimeTZ()
	if err != nil {
		return fmt.Errorf("could not get runtime timezone: %w", err)
	}

	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		return fmt.Errorf("could not get timezone location: %w", err)
	}

	time.Local = location

	logger.WithField("timezoneName", timezoneName).Debug("Runtime timezone set")

	return nil
}

func ParseTimestamp(src string) (time.Time, error) {
	timestamp, err := time.Parse(ResponseTimeFormat, src)

	if err == nil {
		return timestamp, nil
	}

	timestamp, err2 := time.Parse(ResponseTimeFormatAlt, src)

	if err2 == nil {
		return timestamp, nil
	}

	return time.Time{}, fmt.Errorf("could not parse timestamp 2 times: '%w' and '%w'", err, err2)
}
