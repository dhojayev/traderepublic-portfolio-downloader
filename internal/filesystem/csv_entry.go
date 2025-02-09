package filesystem

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
)

type DepotTransactionCSVEntry struct {
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
	TaxAmount      float64 `csv:"Tax amount"`
	InvestedAmount float64 `csv:"-"`
	Documents      []string
}

func NewDepotTransactionCSVEntry(
	id, status, transactionType, assetType, name, instrument string,
	shares, rate, yield, profit, commission, debit, credit, taxAmount, investedAmount float64,
	timestamp internal.DateTime,
	documents []string,
) DepotTransactionCSVEntry {
	return DepotTransactionCSVEntry{
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
		TaxAmount:      taxAmount,
		InvestedAmount: investedAmount,
		Documents:      documents,
	}
}

type CardTransactionCSVEntry struct {
	ID                    string
	Status                string
	Timestamp             internal.DateTime
	Type                  string
	Card                  string
	Merchant              string
	OriginalCurrency      string  `csv:"OriginalCurrency"`
	OriginalAmount        float64 `csv:"Original amount"`
	ExchangeRate          float64
	ReferenceExchangeRate float64 `csv:"Reference exchange rate"`
	Debit                 float64
	Credit                float64
}

func NewCardTransactionCSVEntry(
	id, status, transactionType, card, merchant, originalCurrency string,
	timestamp internal.DateTime,
	originalAmount, exchangeRate, refExchangeRate, debit, credit float64,
) CardTransactionCSVEntry {
	return CardTransactionCSVEntry{
		ID:                    id,
		Status:                status,
		Timestamp:             timestamp,
		Type:                  transactionType,
		Card:                  card,
		Merchant:              merchant,
		OriginalCurrency:      originalCurrency,
		OriginalAmount:        originalAmount,
		ExchangeRate:          exchangeRate,
		ReferenceExchangeRate: refExchangeRate,
		Debit:                 debit,
		Credit:                credit,
	}
}
