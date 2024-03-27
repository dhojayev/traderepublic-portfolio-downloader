package transaction

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
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

func CreateDocumentsFromResponse(resp details.Response) ([]Document, error) {
	documents := make([]Document, 0)

	documentsSection, err := resp.DocumentsSection()
	if err != nil {
		return documents, fmt.Errorf("could not get documents for id %s: %w", resp.ID, err)
	}

	for _, document := range documentsSection.Data {
		url, ok := document.Action.Payload.(string)
		if !ok {
			continue
		}

		documents = append(documents, Document{
			ID:    document.ID,
			URL:   url,
			Date:  document.Detail,
			Title: document.Title,
		})
	}

	return documents, nil
}
