package activity

import (
	"errors"
	"fmt"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	log "github.com/sirupsen/logrus"
)

type ProcessorInterface interface {
	Process(response details.Response) error
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

func (p Processor) Process(response details.Response) error {
	logFields := log.Fields{
		"id": response.ID,
	}

	entryTimestamp, err := p.extractTimestamp(response)
	if err != nil {
		return err
	}

	documents, err := p.builder.Build(response.ID, entryTimestamp, response)
	if err != nil {
		if errors.Is(err, details.ErrSectionTypeNotFound) {
			p.logger.Warnf("document model builder errors: %s", err)

			return nil
		}

		return fmt.Errorf("document model builder error: %w", err)
	}

	for _, doc := range documents {
		err = p.docDownloader.Download(internal.DocumentsBaseDir, doc)
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

func (p Processor) extractTimestamp(response details.Response) (time.Time, error) {
	header, err := response.SectionTypeHeader()
	if err != nil {
		return time.Time{}, fmt.Errorf("could not get header section: %w", err)
	}

	timestamp, err := time.Parse(details.ResponseTimeFormat, header.Data.Timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse header section timestamp: %w", err)
	}

	return timestamp, nil
}
