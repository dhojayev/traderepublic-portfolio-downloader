package document

import "time"

type Model struct {
	TransactionUUID string `gorm:"foreignKey:UUID"`
	ID              string `gorm:"primaryKey"`
	URL             string `gorm:"-"`
	Detail          string
	Title           string
	Timestamp       time.Time
	Filename        string
}

func NewModel(transactionUUID, id, url, detail, title, filename string, timestamp time.Time) Model {
	return Model{
		TransactionUUID: transactionUUID,
		ID:              id,
		URL:             url,
		Detail:          detail,
		Title:           title,
		Timestamp:       timestamp,
		Filename:        filename,
	}
}

func (Model) TableName() string {
	return "documents"
}
