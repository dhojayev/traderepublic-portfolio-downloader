package transaction

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	TypePurchase       = "Purchase"
	TypeSale           = "Sale"
	TypeDividendPayout = "Dividends"
	TypeRoundUp        = "Round up"
	TypeSaveback       = "Saveback"

	InstrumentTypeStocks         = "Stocks"
	InstrumentTypeETF            = "ETF"
	InstrumentTypeCryptocurrency = "Cryptocurrency"
	InstrumentTypeLending        = "Lending"
	InstrumentTypeOther          = "Other"

	isinPrefixLending = "XS"
	isinPrefixCrypto  = "XF000"
	isinSuffixDist    = "(Dist)"
	isinSuffixAcc     = "(Acc)"
)

type Transaction struct {
	UUID string `gorm:"primaryKey"`

	InstrumentID int
	Instrument   Instrument
	Documents    []Document `gorm:"-"`

	Type       string    `gorm:"index"`
	Timestamp  time.Time `gorm:"index"`
	Status     string
	Yield      float64
	Profit     float64
	Shares     float64
	Rate       float64
	Commission float64
	Total      float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func NewTransaction(
	uuid, transactionType, status string,
	yield, profit, shares, rate, commission, total float64,
	timestamp time.Time,
	instrument Instrument,
	documents []Document,
) Transaction {
	return Transaction{
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
		Instrument: instrument,
		Documents:  documents,
	}
}

type Instrument struct {
	gorm.Model

	ISIN string `gorm:"index"`
	Name string
}

func (i Instrument) Type() string {
	instrumentType := InstrumentTypeOther

	switch {
	case strings.HasSuffix(i.Name, isinSuffixDist), strings.HasSuffix(i.Name, isinSuffixAcc):
		instrumentType = InstrumentTypeETF
	case strings.HasPrefix(i.ISIN, isinPrefixCrypto):
		instrumentType = InstrumentTypeCryptocurrency
	case strings.HasPrefix(i.ISIN, isinPrefixLending):
		instrumentType = InstrumentTypeLending
	}

	return instrumentType
}

func NewInstrument(isin, name string) Instrument {
	return Instrument{
		ISIN: isin,
		Name: name,
	}
}
