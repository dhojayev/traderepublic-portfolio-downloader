package document

type Model struct {
	TransactionUUID string `gorm:"foreignKey:UUID"`
	ID              string `gorm:"primaryKey"`
	URL             string `gorm:"-"`
	Detail          string
	Title           string
	Filepath        string
}

func NewModel(transactionUUID, id, url, detail, title, filename string) Model {
	return Model{
		TransactionUUID: transactionUUID,
		ID:              id,
		URL:             url,
		Detail:          detail,
		Title:           title,
		Filepath:        filename,
	}
}

func (Model) TableName() string {
	return "documents"
}
