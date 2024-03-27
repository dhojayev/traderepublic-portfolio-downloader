package transaction

import (
	"errors"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"

	log "github.com/sirupsen/logrus"
)

const csvFilename = "transactions.csv"

type Processor struct {
	deserializer DetailsDeserializerInterface
	factory      filesystem.FactoryInterface
	csvReader    filesystem.CSVReader
	csvWriter    filesystem.CSVWriter
	logger       *log.Logger
}

func NewProcessor(
	deserializer DetailsDeserializerInterface,
	factory filesystem.FactoryInterface,
	csvReader filesystem.CSVReader,
	csvWriter filesystem.CSVWriter,
	logger *log.Logger,
) Processor {
	return Processor{
		deserializer: deserializer,
		factory:      factory,
		csvReader:    csvReader,
		csvWriter:    csvWriter,
		logger:       logger,
	}
}

func (p Processor) Process(response details.Response) error {
	entries, err := p.csvReader.Read(csvFilename)
	if err != nil {
		return fmt.Errorf("csv reader read error: %w", err)
	}

	for _, entry := range entries {
		if entry.ID == response.ID {
			return nil
		}
	}

	transaction, err := p.deserializer.Deserialize(response)
	if err != nil {
		if errors.Is(err, ErrUnsupportedResponse) {
			p.logger.WithField("id", response.ID).Debug(err)

			return ErrUnsupportedResponse
		}

		return fmt.Errorf("deserializer error: %w", err)
	}

	p.logger.WithField("transaction", transaction).Trace("supported transaction detected")

	entry, err := p.factory.Make(transaction)
	if err != nil {
		return fmt.Errorf("could not make csv entry: %w", err)
	}

	if err := p.csvWriter.Write(csvFilename, entry); err != nil {
		return fmt.Errorf("could not save transaction to file: %w", err)
	}

	return nil
}
