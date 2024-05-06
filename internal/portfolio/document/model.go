package document

type Model struct {
	ID    string
	URL   string
	Date  string
	Title string
}

func NewModel(id, url, date, title string) Model {
	return Model{
		ID:    id,
		URL:   url,
		Date:  date,
		Title: title,
	}
}
