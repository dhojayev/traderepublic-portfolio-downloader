//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=downloader.go -destination downloader_mock.go -package=document

package document

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	DownloaderTimeFormat = "2006-01"
	permDir              = 0o700
)

type DownloaderInterface interface {
	Download(baseDir string, document Model) error
}

type Downloader struct {
	logger *log.Logger
}

func NewDownloader(logger *log.Logger) Downloader {
	return Downloader{logger: logger}
}

func (d Downloader) Download(baseDir string, document Model) error {
	dest := baseDir + "/" + document.Filepath
	dir := filepath.Dir(dest)

	if err := os.MkdirAll(dir, permDir); err != nil {
		return fmt.Errorf("could not create directory for document: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, document.URL, nil)
	if err != nil {
		return fmt.Errorf("could not create request for document download: %w", err)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(context.Background()))
	if err != nil {
		return fmt.Errorf("could not download document: %w", err)
	}

	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create document file: %w", err)
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not write document file: %w", err)
	}

	return nil
}
