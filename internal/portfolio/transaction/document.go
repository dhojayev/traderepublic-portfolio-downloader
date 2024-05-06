package transaction

import (
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

type Document struct {
	ID    string
	URL   string
	Date  string
	Title string
}

func NewDocument(id, url, date, title string) Document {
	return Document{
		ID:    id,
		URL:   url,
		Date:  date,
		Title: title,
	}
}

type DocumentDownloader struct {
	logger *log.Logger
}

func NewDocumentDownloader(logger *log.Logger) DocumentDownloader {
	return DocumentDownloader{logger: logger}
}

func (d DocumentDownloader) Download(destDir string, document Document) (string, error) {
	resp, err := grab.Get(destDir, document.URL)
	if err != nil {
		return "", fmt.Errorf("could not download document: %w", err)
	}

	return resp.Filename, nil
}
