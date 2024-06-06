//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=csv_writer.go -destination csv_writer_mock.go -package=filesystem

package filesystem

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

type CSVWriterInterface interface {
	Write(filepath string, entry CSVEntry) error
}

type CSVWriter struct {
	logger *log.Logger
}

func NewCSVWriter(logger *log.Logger) CSVWriter {
	return CSVWriter{
		logger: logger,
	}
}

func (w CSVWriter) Write(filepath string, entry CSVEntry) error {
	entries := []CSVEntry{entry}
	newFile := false

	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		newFile = true
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, permFile)
	if err != nil {
		return fmt.Errorf("could not open csv file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			w.logger.Errorf("csv writer could not close file: %s", err)
		}
	}()

	if newFile {
		err = gocsv.MarshalFile(&entries, file)
		if err != nil {
			return fmt.Errorf("could not write csv file: %w", err)
		}

		return nil
	}

	err = gocsv.MarshalWithoutHeaders(&entries, file)
	if err != nil {
		return fmt.Errorf("could not write csv file: %w", err)
	}

	return nil
}
