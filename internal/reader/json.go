package reader

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type JSONReader struct {
	logger  *log.Logger
	baseDir string
	cursors map[string]uint
}

func NewJSONReader(baseDir string, logger *log.Logger) *JSONReader {
	return &JSONReader{
		baseDir: baseDir,
		logger:  logger,
		cursors: map[string]uint{},
	}
}

func (r *JSONReader) Read(dataType string, req Request) (JSONResponse, error) {
	id, found := req["id"]
	if !found {
		cursor, found := r.cursors[dataType]
		if !found {
			cursor = 1
		}

		r.cursors[dataType] = cursor + 1

		return r.read(fmt.Sprintf("%s/%s/page-%d.json", r.baseDir, dataType, cursor))
	}

	filepath := fmt.Sprintf("%s/%s/%s.json", r.baseDir, dataType, id)

	return r.read(filepath)
}

func (r *JSONReader) read(filepath string) (JSONResponse, error) {
	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		return NewJSONResponse(nil), fmt.Errorf("could not read filepath: %w", err)
	}

	r.logger.
		WithFields(log.Fields{
			"filepath": filepath,
		}).
		Debug("read file contents")

	return NewJSONResponse(fileContents), nil
}
