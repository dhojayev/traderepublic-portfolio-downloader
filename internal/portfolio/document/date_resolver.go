package document

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const ResolverTimeFormat = "02.01.2006"

type DateResolverInterface interface {
	Resolve(transactionTimestamp time.Time, documentDate string) (time.Time, error)
}

type DateResolver struct {
	logger *log.Logger
}

func NewDateResolver(logger *log.Logger) DateResolver {
	return DateResolver{
		logger: logger,
	}
}

func (r DateResolver) Resolve(transactionTimestamp time.Time, documentDate string) (time.Time, error) {
	if documentDate == "" {
		return transactionTimestamp, nil
	}

	date, err := time.Parse(ResolverTimeFormat, documentDate)
	if err != nil {
		return date, fmt.Errorf("could not parse date from document detail: %w", err)
	}

	return date, nil
}
