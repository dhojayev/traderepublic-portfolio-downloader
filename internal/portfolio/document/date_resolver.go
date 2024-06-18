package document

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const ResolverTimeFormat = "02.01.2006"

type DateResolverInterface interface {
	Resolve(parentTimestamp time.Time, documentDate string) time.Time
}

type DateResolver struct {
	logger *log.Logger
}

func NewDateResolver(logger *log.Logger) DateResolver {
	return DateResolver{
		logger: logger,
	}
}

func (r DateResolver) Resolve(parentTimestamp time.Time, documentDate string) time.Time {
	date, err := time.Parse(ResolverTimeFormat, documentDate)
	if err != nil {
		return parentTimestamp
	}

	return date
}
