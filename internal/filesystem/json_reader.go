package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"

	log "github.com/sirupsen/logrus"
)

type JSONReader struct {
	logger  *log.Logger
	baseDir string
	cursors map[string]uint
}

func NewJSONReader(baseDir string, logger *log.Logger) JSONReader {
	return JSONReader{
		baseDir: baseDir,
		logger:  logger,
		cursors: map[string]uint{},
	}
}

func (r JSONReader) Read(dataType string, data map[string]any) (portfolio.OutputDataInterface, error) { //nolint:ireturn
	id, found := data["id"]
	if !found {
		cursor, found := r.cursors[dataType]
		if !found {
			cursor = 0
			r.cursors[dataType] = cursor
		}

		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", r.baseDir, dataType))
		if err != nil {
			return OutputData{}, fmt.Errorf("could not read dir: %w", err)
		}

		var filteredEntries []fs.DirEntry

		for _, entry := range entries {
			if entry.IsDir() || !strings.Contains(entry.Name(), ".json") {
				continue
			}

			filteredEntries = append(filteredEntries, entry)
		}

		entry := filteredEntries[cursor]
		filepath := fmt.Sprintf("%s/%s/%s", r.baseDir, dataType, entry.Name())
		r.cursors[dataType]++

		return r.read(filepath)
	}

	filepath := fmt.Sprintf("%s/%s/%s.json", r.baseDir, dataType, id)

	return r.read(filepath)
}

func (r JSONReader) read(filepath string) (OutputData, error) {
	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		return OutputData{}, fmt.Errorf("could not read filepath: %w", err)
	}

	r.logger.
		WithFields(log.Fields{
			"filepath": filepath,
			"contents": string(fileContents),
		}).
		Trace("read file contents")

	return OutputData{data: fileContents}, nil
}
