package filesystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	baseDir  = "responses"
	permDir  = 0o700
	permFile = 0o600
)

type JSONWriter struct {
	logger  *log.Logger
	cursors map[string]uint
}

func NewJSONWriter(logger *log.Logger) *JSONWriter {
	return &JSONWriter{
		logger:  logger,
		cursors: map[string]uint{},
	}
}

func (w *JSONWriter) Bytes(dir string, data []byte) error {
	if dir == "" {
		return errors.New("writer: dir cannot be empty")
	}

	var dataMap map[string]any

	if err := json.Unmarshal(data, &dataMap); err != nil {
		return fmt.Errorf("could not unmarshal data bytes to write to file: %w", err)
	}

	filename, err := w.GenerateFilename(dir, dataMap)
	if err != nil {
		return err
	}

	return w.write(dir, filename, dataMap)
}

func (w *JSONWriter) GenerateFilename(dir string, dataMap map[string]any) (string, error) {
	idVal, found := dataMap["id"]

	if found {
		id, ok := idVal.(string)
		if !ok {
			return "", errors.New("could not convert id into string")
		}

		return id, nil
	}

	cursor, found := w.cursors[dir]
	if !found {
		cursor = 1
	}

	w.cursors[dir] = cursor + 1

	return fmt.Sprintf("page-%d", cursor), nil
}

func (w *JSONWriter) write(dir, filename string, dataMap map[string]any) error {
	data, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal data bytes to write to file: %w", err)
	}

	destDir := baseDir + "/" + dir
	filepath := fmt.Sprintf("./%s/%s.json", destDir, filename)

	if err := w.createDir(destDir); err != nil {
		return err
	}

	if err := os.WriteFile(filepath, data, permFile); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	w.logger.
		WithFields(log.Fields{
			"filepath": filepath,
			"contents": string(data),
		}).
		Trace("wrote file")

	return nil
}

func (w *JSONWriter) createDir(destDir string) error {
	_, err := os.Stat(destDir)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("could not stat destination dir: %w", err)
	}

	if err := os.MkdirAll(destDir, permDir); err != nil {
		return fmt.Errorf("could not create destination dir: %w", err)
	}

	return nil
}
