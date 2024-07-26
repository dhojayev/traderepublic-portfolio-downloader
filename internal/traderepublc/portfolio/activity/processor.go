package activity

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
)

type ProcessorInterface interface {
	Process(response details.NormalizedResponse) error
}

type Processor struct {
	builder       document.ModelBuilderInterface
	docDownloader document.DownloaderInterface
	logger        *log.Logger
}

func NewProcessor(
	builder document.ModelBuilderInterface,
	docDownloader document.DownloaderInterface,
	logger *log.Logger,
) Processor {
	return Processor{
		builder:       builder,
		docDownloader: docDownloader,
		logger:        logger,
	}
}

func (p Processor) Process(response details.NormalizedResponse) error {
	logFields := log.Fields{
		"id": response.ID,
	}

	entryTimestamp, err := p.extractTimestamp(response)
	if err != nil {
		return err
	}

	documents, err := p.builder.Build(response.ID, entryTimestamp, response)
	if err != nil {
		return fmt.Errorf("document model builder error: %w", err)
	}

	for _, doc := range documents {
		err = p.docDownloader.Download(internal.ActivityLogDocumentsBaseDir, doc)
		if err == nil {
			continue
		}

		if errors.Is(err, document.ErrDocumentExists) {
			continue
		}

		p.logger.WithFields(logFields).Warnf("Document downloader error: %s", err)
	}

	return nil
}

func (p Processor) extractTimestamp(response details.NormalizedResponse) (time.Time, error) {
	timestamp, err := time.Parse(details.ResponseTimeFormat, response.Header.Data.Timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse header section timestamp: %w", err)
	}

	return timestamp, nil
}
