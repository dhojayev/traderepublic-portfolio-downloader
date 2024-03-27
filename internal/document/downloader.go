package document

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
)

type Downloader struct{}

func NewDownloader() Downloader {
	return Downloader{}
}

func (d Downloader) Download(resp details.Response) ([]string, error) {
	filepaths := make([]string, 0)

	entries, err := NewDownloadEntriesFromResponse(resp)
	if err != nil {
		return filepaths, err
	}

	for _, entry := range entries {
		filepath, err := d.download(entry)
		if err != nil {
			return filepaths, err
		}

		filepaths = append(filepaths, filepath)
	}

	return filepaths, nil
}

func (d Downloader) download(_ DownloadEntry) (string, error) {
	return "", nil
}
