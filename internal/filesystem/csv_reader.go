package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

type CSVReader struct {
	logger *log.Logger
}

func NewCSVReader(logger *log.Logger) CSVReader {
	return CSVReader{
		logger: logger,
	}
}

func (r CSVReader) Read(filepath string) ([]CSVEntry, error) {
	var entries []CSVEntry

	file, err := os.OpenFile(filepath, os.O_RDWR, os.ModePerm)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return entries, nil
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			r.logger.Errorf("csv reader could not close file: %s", err)
		}
	}()

	if err := gocsv.UnmarshalFile(file, &entries); err != nil { // Load clients from file
		return entries, fmt.Errorf("csv unmarshall error: %w", err)
	}

	return entries, nil
}
