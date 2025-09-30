package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
)

const (
	filePermissions = 0600
)

type ResponseWriter struct {
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

func (w *ResponseWriter) Bytes(filename string, data []byte) error {
	formattedFilename := filepath.Join(internal.ResponseBaseDir, filename+".json")

	err := os.MkdirAll(filepath.Dir(formattedFilename), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	err = os.WriteFile(formattedFilename, data, filePermissions)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
