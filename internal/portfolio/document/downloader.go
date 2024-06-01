//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=downloader.go -destination downloader_mock.go -package=document

package document

import (
	"fmt"
	"os"

	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

const (
	DownloaderTimeFormat = "2006-01-02"
	permDir              = 0o700
)

type DownloaderInterface interface {
	Download(baseDir string, document Model) (string, error)
}

type Downloader struct {
	logger *log.Logger
}

func NewDownloader(logger *log.Logger) Downloader {
	return Downloader{logger: logger}
}

func (d Downloader) Download(baseDir string, document Model) (string, error) {
	destDir := fmt.Sprintf("%s/%s/", baseDir, document.Timestamp.Format(DownloaderTimeFormat))

	if err := os.MkdirAll(destDir, permDir); err != nil {
		return "", fmt.Errorf("could not create download base dir: %w", err)
	}

	resp, err := grab.Get(destDir, document.URL)
	if err != nil {
		return "", fmt.Errorf("could not download document: %w", err)
	}

	return resp.Filename, nil
}
