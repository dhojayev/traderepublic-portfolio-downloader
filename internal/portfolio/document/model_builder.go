package document

import (
	"fmt"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	log "github.com/sirupsen/logrus"
)

type ModelBuilderInterface interface {
	Build(transactionUUID string, transactionTimestamp time.Time, response details.Response) ([]Model, error)
}

type ModelBuilder struct {
	dateResolver DateResolverInterface
	logger       *log.Logger
}

func NewModelBuilder(dateResolver DateResolverInterface, logger *log.Logger) ModelBuilder {
	return ModelBuilder{
		dateResolver: dateResolver,
		logger:       logger,
	}
}

func (b ModelBuilder) Build(
	_ string,
	transactionTimestamp time.Time,
	response details.Response,
) ([]Model, error) {
	documents := make([]Model, 0)

	documentsSection, err := response.SectionTypeDocuments()
	if err != nil {
		return documents, fmt.Errorf("could not get documents section: %w", err)
	}

	for _, doc := range documentsSection.Data {
		url, ok := doc.Action.Payload.(string)
		if !ok {
			continue
		}

		documentDate, err := b.dateResolver.Resolve(transactionTimestamp, doc.Detail)
		if err != nil {
			return documents, fmt.Errorf("document date resolver errors: %w", err)
		}

		documents = append(documents, NewModel("", doc.ID, url, doc.Detail, doc.Title, "", documentDate))
	}

	return documents, nil
}
