package transaction

import (
	"fmt"
	"strings"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
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

	InstrumentTypeStocks         = "Stocks"
	InstrumentTypeETF            = "ETF"
	InstrumentTypeCryptocurrency = "Cryptocurrency"
	InstrumentTypeLending        = "Lending"
	InstrumentTypeCash           = "Cash"
	InstrumentTypeOther          = "Other"

	isinPrefixLending = "XS"
	isinPrefixCrypto  = "XF000"
	isinSuffixDist    = "(Dist)"
	isinSuffixAcc     = "(Acc)"
)

type Model struct {
	UUID string `gorm:"primaryKey"`

	InstrumentISIN *string
	Instrument     Instrument
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
	instrument Instrument,
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

type Instrument struct {
	ISIN string `gorm:"primaryKey"`
	Name string
	Icon string
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

func (i Instrument) IconURL() string {
	return fmt.Sprintf("https://assets.traderepublic.com/img/%s/light.min.svg", i.Icon)
}

func NewInstrument(isin, name, icon string) Instrument {
	return Instrument{
		ISIN: isin,
		Name: name,
		Icon: icon,
	}
}
