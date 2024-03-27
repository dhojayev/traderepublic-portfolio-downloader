package document

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
)

type DownloadEntry struct {
	src details.ResponseSectionTypeDocumentsData
}

func NewDownloadEntry(documentSection details.ResponseSectionTypeDocumentsData) DownloadEntry {
	return DownloadEntry{
		src: documentSection,
	}
}

func NewDownloadEntriesFromResponse(resp details.Response) ([]DownloadEntry, error) {
	entries := make([]DownloadEntry, 0)

	section, err := resp.DocumentsSection()
	if err != nil {
		return entries, fmt.Errorf("could not get document section from response: %w", err)
	}

	for _, data := range section.Data {
		entries = append(entries, NewDownloadEntry(data))
	}

	return entries, nil
}

func (e DownloadEntry) Filename() string {
	return ""
}

func (e DownloadEntry) Dir() string {
	return ""
}
