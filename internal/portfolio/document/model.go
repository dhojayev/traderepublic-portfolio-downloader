package document

import "time"

type Model struct {
	ID        string
	URL       string
	Detail    string
	Title     string
	Timestamp time.Time
	Filepath  string
}

func NewModel(id, url, detail, title, filepath string, timestamp time.Time) Model {
	return Model{
		ID:        id,
		URL:       url,
		Detail:    detail,
		Title:     title,
		Timestamp: timestamp,
		Filepath:  filepath,
	}
}
