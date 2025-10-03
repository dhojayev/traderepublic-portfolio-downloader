package transaction

import (
	"fmt"
	"time"
)

type CSVEntry struct {
	ID             string
	Status         string
	Timestamp      DateTime
	Type           TransactionType
	AssetType      string `csv:"Asset type"`
	Name           string
	Instrument     string
	Shares         float64
	Rate           float64 `csv:"Realized yield"`
	Yield          float64 `csv:"Realized PnL"`
	Profit         float64 `csv:"Realized PnL"`
	Commission     float64
	Debit          float64
	Credit         float64
	TaxAmount      float64 `csv:"Tax amount"`
	InvestedAmount float64 `csv:"-"`
	Documents      []string
}

type TransactionType string

const (
	Buy      TransactionType = "buy"
	Sell     TransactionType = "sell"
	Dividend TransactionType = "dividend"
)

func NewCSVEntryBuilder() *CSVEntryBuilder {
	return &CSVEntryBuilder{}
}

type CSVEntryBuilder struct {
	ID             string
	Status         string
	Timestamp      DateTime
	Type           TransactionType
	AssetType      string `csv:"Asset type"`
	Name           string
	Instrument     string
	Shares         float64
	Rate           float64 `csv:"Realized yield"`
	Yield          float64 `csv:"Realized PnL"`
	Profit         float64 `csv:"Realized PnL"`
	Commission     float64
	Debit          float64
	Credit         float64
	TaxAmount      float64 `csv:"Tax amount"`
	InvestedAmount float64 `csv:"-"`
	Documents      []string
}

func (b *CSVEntryBuilder) WithID(id string) *CSVEntryBuilder {
	b.ID = id
	
	return b
}

func (b *CSVEntryBuilder) WithStatus(status string) *CSVEntryBuilder {
	b.Status = status
	
	return b
}

func (b *CSVEntryBuilder) WithTimestamp(timestamp DateTime) *CSVEntryBuilder {
	b.Timestamp = timestamp
	
	return b
}

func (b *CSVEntryBuilder) WithType(transactionType TransactionType) *CSVEntryBuilder {
	b.Type = transactionType
	
	return b
}

func (b *CSVEntryBuilder) WithAssetType(assetType string) *CSVEntryBuilder {
	b.AssetType = assetType
	
	return b
}

func (b *CSVEntryBuilder) WithName(name string) *CSVEntryBuilder {
	b.Name = name
	
	return b
}

func (b *CSVEntryBuilder) WithInstrument(instrument string) *CSVEntryBuilder {
	b.Instrument = instrument
	
	return b
}

func (b *CSVEntryBuilder) WithShares(shares float64) *CSVEntryBuilder {
	b.Shares = shares
	
	return b
}

func (b *CSVEntryBuilder) WithRate(rate float64) *CSVEntryBuilder {
	b.Rate = rate
	
	return b
}

func (b *CSVEntryBuilder) WithYield(yield float64) *CSVEntryBuilder {
	b.Yield = yield
	
	return b
}

func (b *CSVEntryBuilder) WithProfit(profit float64) *CSVEntryBuilder {
	b.Profit = profit
	
	return b
}

func (b *CSVEntryBuilder) WithCommission(commission float64) *CSVEntryBuilder {
	b.Commission = commission
	
	return b
}

func (b *CSVEntryBuilder) WithDebit(debit float64) *CSVEntryBuilder {
	b.Debit = debit
	
	return b
}

func (b *CSVEntryBuilder) WithCredit(credit float64) *CSVEntryBuilder {
	b.Credit = credit
	
	return b
}

func (b *CSVEntryBuilder) WithTaxAmount(taxAmount float64) *CSVEntryBuilder {
	b.TaxAmount = taxAmount
	
	return b
}

func (b *CSVEntryBuilder) WithInvestedAmount(investedAmount float64) *CSVEntryBuilder {
	b.InvestedAmount = investedAmount
	
	return b
}

func (b *CSVEntryBuilder) AddDocument(document string) *CSVEntryBuilder {
	b.Documents = append(b.Documents, document)
	
	return b
}

func (b *CSVEntryBuilder) Build() CSVEntry {
	return CSVEntry{
		ID:             b.ID,
		Status:         b.Status,
		Timestamp:      b.Timestamp,
		Type:           b.Type,
		AssetType:      b.AssetType,
		Name:           b.Name,
		Instrument:     b.Instrument,
		Shares:         b.Shares,
		Rate:           b.Rate,
		Yield:          b.Yield,
		Profit:         b.Profit,
		Commission:     b.Commission,
		Debit:          b.Debit,
		Credit:         b.Credit,
		TaxAmount:      b.TaxAmount,
		InvestedAmount: b.InvestedAmount,
		Documents:      b.Documents,
	}
}

type DateTime struct {
	time.Time
}

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.In(time.Local).Format(time.RFC822Z), nil
}

func (date *DateTime) UnmarshalCSV(csv string) error {
	t, err := time.Parse(time.RFC822Z, csv)
	if err != nil {
		return fmt.Errorf("could not parse datetime: %w", err)
	}

	date.Time = t

	return nil
}
