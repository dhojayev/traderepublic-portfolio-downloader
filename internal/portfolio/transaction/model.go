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

type Purchase struct {
	gorm.Model

	TransactionID int
	Transaction   Transaction
}

func NewPurchase(transaction Transaction) Purchase {
	transaction.Type = TypePurchase

	return Purchase{
		Transaction: transaction,
	}
}

type Sale struct {
	gorm.Model

	Yield  float64
	Profit float64

	TransactionID int
	Transaction   Transaction
}

func NewSale(
	yield, profit float64,
	transaction Transaction,
) Sale {
	transaction.Type = TypeSale

	return Sale{
		Yield:       yield,
		Profit:      profit,
		Transaction: transaction,
	}
}

type Benefit struct {
	gorm.Model

	TransactionID int
	Transaction   Transaction
}

func NewBenefit(benefitType string, transaction Transaction) Benefit {
	transaction.Type = benefitType

	return Benefit{
		Transaction: transaction,
	}
}

func (b Benefit) IsTypeRoundUp() bool {
	return b.Transaction.Type == TypeRoundUp
}

type DividendPayout struct {
	gorm.Model

	Profit float64

	TransactionID int
	Transaction   Transaction
}

func NewDividendPayout(profit float64, transaction Transaction) DividendPayout {
	transaction.Type = TypeDividendPayout

	return DividendPayout{
		Profit:      profit,
		Transaction: transaction,
	}
}

type Transaction struct {
	gorm.Model

	UUID       string
	Type       string
	Timestamp  time.Time
	Status     string
	Yield      float64
	Profit     float64
	Shares     float64
	Rate       float64
	Commission float64
	Total      float64

	InstrumentID int
	Instrument   Instrument

	Documents []Document `gorm:"-"`
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

	ISIN string
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
