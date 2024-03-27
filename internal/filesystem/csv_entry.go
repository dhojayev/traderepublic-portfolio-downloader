package filesystem

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
)

type CSVEntry struct {
	ID             string
	Status         string
	Timestamp      internal.DateTime
	Type           string
	AssetType      string `csv:"Asset type"`
	Name           string
	Instrument     string
	Shares         float64
	Rate           float64
	Yield          float64 `csv:"Realized yield"`
	Profit         float64 `csv:"Realized PnL"`
	Commission     float64
	Debit          float64
	Credit         float64
	PortfolioValue float64 `csv:"Portfolio value"`
}

func NewCSVEntry(
	id, status, transactionType, assetType, name, instrument string,
	shares, rate, yield, profit, commission, debit, credit, portfolioValue float64,
	timestamp internal.DateTime,
) CSVEntry {
	return CSVEntry{
		ID:             id,
		Status:         status,
		Timestamp:      timestamp,
		Type:           transactionType,
		AssetType:      assetType,
		Name:           name,
		Instrument:     instrument,
		Shares:         shares,
		Rate:           rate,
		Yield:          yield,
		Profit:         profit,
		Commission:     commission,
		Debit:          debit,
		Credit:         credit,
		PortfolioValue: portfolioValue,
	}
}

type FactoryInterface interface {
	Make(valueObject any) (CSVEntry, error)
}
