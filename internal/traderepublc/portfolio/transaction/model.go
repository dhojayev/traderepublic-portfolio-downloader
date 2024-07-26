package transaction

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/instrument"
)

const (
	TypePurchase       = "Purchase"
	TypeSale           = "Sale"
	TypeDividendPayout = "Dividends"
	TypeRoundUp        = "Round up"
	TypeSaveback       = "Saveback"
	TypeDeposit        = "Deposit"
	TypeWithdrawal     = "Withdrawal"
	TypeInterestPayout = "Interest payout"
)

type Model struct {
	UUID string `gorm:"primaryKey"`

	InstrumentISIN *string
	Instrument     instrument.Model
	Documents      []document.Model `gorm:"-"`

	Type       string    `gorm:"index"`
	Timestamp  time.Time `gorm:"index"`
	Status     string
	Yield      float64
	Profit     float64
	Shares     float64
	Rate       float64
	Commission float64
	Total      float64
	TaxAmount  float64
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time `gorm:"index"`
}

func NewTransaction(
	uuid, transactionType, status string,
	yield, profit, shares, rate, commission, total, tax float64,
	timestamp time.Time,
	instrument instrument.Model,
	documents []document.Model,
) Model {
	return Model{
		UUID:       uuid,
		Type:       transactionType,
		Timestamp:  timestamp,
		Status:     status,
		Yield:      yield,
		Profit:     profit,
		Shares:     shares,
		Rate:       rate,
		Commission: commission,
		Total:      total,
		TaxAmount:  tax,
		Instrument: instrument,
		Documents:  documents,
	}
}

func (Model) TableName() string {
	return "transactions"
}
