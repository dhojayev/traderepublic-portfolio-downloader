package transaction

import (
	"errors"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"

	log "github.com/sirupsen/logrus"
)

const (
	csvFilename     = "./transactions.csv"
	documentBaseDir = "./documents"
)

type Processor struct {
	builderFactory     ModelBuilderFactoryInterface
	transactionRepo    RepositoryInterface
	factory            CSVEntryFactory
	csvReader          filesystem.CSVReader
	csvWriter          filesystem.CSVWriter
	documentDownloader document.DownloaderInterface
	logger             *log.Logger
}

func NewProcessor(
	builderFactory ModelBuilderFactoryInterface,
	transactionRepo RepositoryInterface,
	factory CSVEntryFactory,
	csvReader filesystem.CSVReader,
	csvWriter filesystem.CSVWriter,
	documentDownloader document.DownloaderInterface,
	logger *log.Logger,
) Processor {
	return Processor{
		builderFactory:     builderFactory,
		transactionRepo:    transactionRepo,
		factory:            factory,
		csvReader:          csvReader,
		csvWriter:          csvWriter,
		documentDownloader: documentDownloader,
		logger:             logger,
	}
}

//nolint:cyclop
func (p Processor) Process(eventType transactions.EventType, response details.Response) error {
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
		if errors.Is(err, ErrUnsupportedType) {
			p.logger.WithFields(logFields).Debugf("builder factory error: %s", err)

			return ErrUnsupportedType
		}

		return fmt.Errorf("builder factory error: %w", err)
	}

	transaction, err := builder.Build()
	if err != nil {
		return fmt.Errorf("builder error: %w", err)
	}

	for i, document := range transaction.Documents {
		filename, err := p.documentDownloader.Download(documentBaseDir, document)
		if err != nil {
			p.logger.WithFields(logFields).Warnf("document downloader error: %s", err)

			continue
		}

		document.Filename = filename
		transaction.Documents[i] = document
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

	return nil
}
