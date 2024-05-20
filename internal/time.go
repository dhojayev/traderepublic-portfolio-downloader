package internal

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thlib/go-timezone-local/tzlocal"

	// load timezone data.
	_ "time/tzdata"
)

const DefaultTimeFormat = "2006-01-02T15:04:05-0700"

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
