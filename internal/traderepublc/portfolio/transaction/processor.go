//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=processor.go -destination processor_mock.go -package=transaction

package transaction

import (
	"errors"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"

	log "github.com/sirupsen/logrus"
)

const (
	csvFilename     = "./transactions.csv"
	documentBaseDir = "./documents/transactions"
)

type ProcessorInterface interface {
	Process(eventType transactions.EventType, response details.NormalizedResponse) error
}

type Processor struct {
	builderFactory  ModelBuilderFactoryInterface
	transactionRepo RepositoryInterface
	factory         CSVEntryFactory
	csvReader       filesystem.CSVReaderInterface
	csvWriter       filesystem.CSVWriterInterface
	docDownloader   document.DownloaderInterface
	logger          *log.Logger
}

func NewProcessor(
	builderFactory ModelBuilderFactoryInterface,
	transactionRepo RepositoryInterface,
	factory CSVEntryFactory,
	csvReader filesystem.CSVReaderInterface,
	csvWriter filesystem.CSVWriterInterface,
	docDownloader document.DownloaderInterface,
	logger *log.Logger,
) Processor {
	return Processor{
		builderFactory:  builderFactory,
		transactionRepo: transactionRepo,
		factory:         factory,
		csvReader:       csvReader,
		csvWriter:       csvWriter,
		docDownloader:   docDownloader,
		logger:          logger,
	}
}

//nolint:cyclop
func (p Processor) Process(eventType transactions.EventType, response details.NormalizedResponse) error {
	csvEntries, err := p.csvReader.Read(csvFilename)
	if err != nil {
		return fmt.Errorf("csv reader read error: %w", err)
	}

	for _, entry := range csvEntries {
		if entry.ID == response.ID {
			return nil
		}
	}

	logFields := log.Fields{
		"id": response.ID,
	}

	builder, err := p.builderFactory.Create(eventType, response)
	if err != nil {
		if errors.Is(err, ErrModelBuilderUnsupportedType) {
			p.logger.WithFields(logFields).Debugf("builder factory error: %s", err)

			return ErrModelBuilderUnsupportedType
		}

		return fmt.Errorf("builder factory error: %w", err)
	}

	transaction, err := builder.Build()
	if err != nil {
		return fmt.Errorf("builder error: %w", err)
	}

	if err := p.transactionRepo.Create(&transaction); err != nil {
		return fmt.Errorf("could not create on repo: %w", err)
	}

	entry, err := p.factory.Make(transaction)
	if err != nil {
		return fmt.Errorf("could not make csv entry: %w", err)
	}

	if err := p.csvWriter.Write(csvFilename, entry); err != nil {
		return fmt.Errorf("could not save transaction to file: %w", err)
	}

	for _, doc := range transaction.Documents {
		err = p.docDownloader.Download(documentBaseDir, doc)
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
