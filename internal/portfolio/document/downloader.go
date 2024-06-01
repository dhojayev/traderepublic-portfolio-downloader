//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=downloader.go -destination downloader_mock.go -package=document

package document

import (
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

const (
	DownloaderTimeFormat = "2006-01"
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
	dest := baseDir + "/" + document.Filepath

	resp, err := grab.Get(dest, document.URL)
	if err != nil {
		return "", fmt.Errorf("could not download document: %w", err)
	}

	return resp.Filename, nil
}
