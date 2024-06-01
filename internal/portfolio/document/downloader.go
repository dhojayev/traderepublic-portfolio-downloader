//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=downloader.go -destination downloader_mock.go -package=document

package document

import (
	"fmt"
	"os"

	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

const permDir = 0o700

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
	if err := os.Mkdir(destDir, permDir); err != nil {
		return "", fmt.Errorf("could not create download dir: %w", err)
	}

	resp, err := grab.Get(destDir, document.URL)
	if err != nil {
		return "", fmt.Errorf("could not download document: %w", err)
	}

	return resp.Filename, nil
}
