package document

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
)

type ModelBuilderInterface interface {
	Build(transactionUUID string, parentTimestamp time.Time, response details.NormalizedResponse) ([]Model, error)
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
	parentUUID string,
	parentTimestamp time.Time,
	response details.NormalizedResponse,
) ([]Model, error) {
	documents := make([]Model, 0)

	for _, doc := range response.Documents.Data {
		url, ok := doc.Action.Payload.(string)
		if !ok {
			continue
		}

		documentDate := b.dateResolver.Resolve(parentTimestamp, doc.Detail)
		filepath := fmt.Sprintf("%s/%s/%s.pdf", documentDate.Format(DownloaderTimeFormat), parentUUID, doc.Title)
		documents = append(documents, NewModel(parentUUID, doc.ID, url, doc.Detail, doc.Title, filepath))
	}

	if len(documents) == 0 {
		return nil, nil
	}

	return documents, nil
}
