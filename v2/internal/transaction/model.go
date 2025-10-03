package transaction

type Model struct {
	ID             string
	Status         string
	Timestamp      CSVDateTime
	Type           TransactionType
	AssetType      string `csv:"Asset type"`
	AssetName      string
	ISIN           string
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
	TypePurchase       TransactionType = "Purchase"
	TypeSale           TransactionType = "Sale"
	TypeDividendPayout TransactionType = "Dividends"
	TypeRoundUp        TransactionType = "Round up"
	TypeSaveback       TransactionType = "Saveback"
	TypeDeposit        TransactionType = "Deposit"
	TypeWithdrawal     TransactionType = "Withdrawal"
	TypeInterestPayout TransactionType = "Interest payout"
)

func NewModelBuilder() *ModelBuilder {
	return &ModelBuilder{}
}

type ModelBuilder struct {
	ID             string
	Status         string
	Timestamp      CSVDateTime
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

func (b *ModelBuilder) WithID(id string) *ModelBuilder {
	b.ID = id

	return b
}

func (b *ModelBuilder) WithStatus(status string) *ModelBuilder {
	b.Status = status

	return b
}

func (b *ModelBuilder) WithTimestamp(timestamp CSVDateTime) *ModelBuilder {
	b.Timestamp = timestamp

	return b
}

func (b *ModelBuilder) WithType(transactionType TransactionType) *ModelBuilder {
	b.Type = transactionType

	return b
}

func (b *ModelBuilder) WithAssetType(assetType string) *ModelBuilder {
	b.AssetType = assetType

	return b
}

func (b *ModelBuilder) WithName(name string) *ModelBuilder {
	b.Name = name

	return b
}

func (b *ModelBuilder) WithInstrument(instrument string) *ModelBuilder {
	b.Instrument = instrument

	return b
}

func (b *ModelBuilder) WithShares(shares float64) *ModelBuilder {
	b.Shares = shares

	return b
}

func (b *ModelBuilder) WithRate(rate float64) *ModelBuilder {
	b.Rate = rate

	return b
}

func (b *ModelBuilder) WithYield(yield float64) *ModelBuilder {
	b.Yield = yield

	return b
}

func (b *ModelBuilder) WithProfit(profit float64) *ModelBuilder {
	b.Profit = profit

	return b
}

func (b *ModelBuilder) WithCommission(commission float64) *ModelBuilder {
	b.Commission = commission

	return b
}

func (b *ModelBuilder) WithDebit(debit float64) *ModelBuilder {
	b.Debit = debit

	return b
}

func (b *ModelBuilder) WithCredit(credit float64) *ModelBuilder {
	b.Credit = credit

	return b
}

func (b *ModelBuilder) WithTaxAmount(taxAmount float64) *ModelBuilder {
	b.TaxAmount = taxAmount

	return b
}

func (b *ModelBuilder) WithInvestedAmount(investedAmount float64) *ModelBuilder {
	b.InvestedAmount = investedAmount

	return b
}

func (b *ModelBuilder) AddDocument(document string) *ModelBuilder {
	b.Documents = append(b.Documents, document)

	return b
}

func (b *ModelBuilder) Build() Model {
	return Model{
		ID:             b.ID,
		Status:         b.Status,
		Timestamp:      b.Timestamp,
		Type:           b.Type,
		AssetType:      b.AssetType,
		AssetName:      b.Name,
		ISIN:           b.Instrument,
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
