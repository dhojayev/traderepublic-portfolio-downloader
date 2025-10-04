//go:generate go tool mockgen -source=csv_writer.go -destination csv_writer_mock.go -package=file

package file

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/gocarina/gocsv"
)

const permFile = 0o600

type CSVWriterInterface interface {
	Write(filepath string, entry transaction.Model) error
}

type CSVWriter struct {
	mu *sync.Mutex
}

func NewCSVWriter() *CSVWriter {
	return &CSVWriter{
		mu: &sync.Mutex{},
	}
}

func (w *CSVWriter) Write(filepath string, entry transaction.Model) error {
	w.mu.Lock()

	defer w.mu.Unlock()

	entries := []transaction.Model{entry}
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
			slog.Error("csv writer could not close file", "err", err)
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
