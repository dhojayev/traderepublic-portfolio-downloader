package document

import (
	"fmt"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	log "github.com/sirupsen/logrus"
)

type ModelBuilderInterface interface {
	Build(response details.Response) ([]Model, error)
}

type ModelBuilder struct {
	logger *log.Logger
}

func NewModelBuilder(logger *log.Logger) ModelBuilder {
	return ModelBuilder{
		logger: logger,
	}
}

func (b ModelBuilder) Build(response details.Response) ([]Model, error) {
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

		documents = append(documents, NewModel(doc.ID, url, doc.Detail, doc.Title, "", time.Time{}))
	}

	return documents, nil
}
