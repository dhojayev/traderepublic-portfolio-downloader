//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=downloader.go -destination downloader_mock.go -package=document

package document

import (
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

type DownloaderInterface interface {
	Download(destDir string, document Model) (string, error)
}

type Downloader struct {
	logger *log.Logger
}

func NewDownloader(logger *log.Logger) Downloader {
	return Downloader{logger: logger}
}

func (d Downloader) Download(destDir string, document Model) (string, error) {
	resp, err := grab.Get(destDir, document.URL)
	if err != nil {
		return "", fmt.Errorf("could not download document: %w", err)
	}

	return resp.Filename, nil
}
