package transaction

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
)

type CSVEntryFactory struct {
	logger *log.Logger
}

func NewCSVEntryFactory(logger *log.Logger) CSVEntryFactory {
	return CSVEntryFactory{
		logger: logger,
	}
}

func (f CSVEntryFactory) Make(valueObject any) (filesystem.CSVEntry, error) {
	switch transaction := valueObject.(type) {
	case Purchase:
		return f.FromPurchase(transaction), nil
	case Sale:
		return f.FromSale(transaction), nil
	case Benefit:
		return f.FromBenefit(transaction), nil
	case DividendPayout:
		return f.FromDividendPayout(transaction), nil
	}

	return filesystem.CSVEntry{}, errors.New("unsupported value object received")
}

func (f CSVEntryFactory) FromPurchase(purchase Purchase) filesystem.CSVEntry {
	investedAmount := purchase.MonetaryValues.Total - purchase.MonetaryValues.Commission

	return filesystem.NewCSVEntry(
		purchase.BaseTransaction.UUID,
		purchase.BaseTransaction.Status,
		purchase.BaseTransaction.Type,
		purchase.Asset.Type,
		purchase.Asset.Name,
		purchase.Asset.Instrument,
		purchase.Asset.Shares,
		purchase.MonetaryValues.Rate,
		0,
		0,
		purchase.MonetaryValues.Commission,
		purchase.MonetaryValues.Total,
		0,
		investedAmount,
		internal.DateTime{Time: purchase.BaseTransaction.Timestamp},
	)
}

func (f CSVEntryFactory) FromSale(sale Sale) filesystem.CSVEntry {
	investedAmount := -(sale.MonetaryValues.Total - sale.Profit + sale.Purchase.MonetaryValues.Commission)

	return filesystem.NewCSVEntry(
		sale.BaseTransaction.UUID,
		sale.BaseTransaction.Status,
		sale.BaseTransaction.Type,
		sale.Asset.Type,
		sale.Asset.Name,
		sale.Asset.Instrument,
		-sale.Asset.Shares,
		sale.MonetaryValues.Rate,
		sale.Yield,
		sale.Profit,
		sale.MonetaryValues.Commission,
		0,
		sale.MonetaryValues.Total,
		investedAmount,
		internal.DateTime{Time: sale.BaseTransaction.Timestamp},
	)
}

func (f CSVEntryFactory) FromBenefit(benefit Benefit) filesystem.CSVEntry {
	var debit float64

	if benefit.IsTypeRoundUp() {
		debit = benefit.MonetaryValues.Total
	}

	return filesystem.NewCSVEntry(
		benefit.BaseTransaction.UUID,
		benefit.BaseTransaction.Status,
		benefit.BaseTransaction.Type,
		benefit.Asset.Type,
		benefit.Asset.Name,
		benefit.Asset.Instrument,
		benefit.Asset.Shares,
		benefit.MonetaryValues.Rate,
		0,
		0,
		benefit.MonetaryValues.Commission,
		debit,
		0,
		benefit.MonetaryValues.Total,
		internal.DateTime{Time: benefit.BaseTransaction.Timestamp},
	)
}

func (f CSVEntryFactory) FromDividendPayout(payout DividendPayout) filesystem.CSVEntry {
	return filesystem.NewCSVEntry(
		payout.BaseTransaction.UUID,
		payout.BaseTransaction.Status,
		payout.BaseTransaction.Type,
		payout.Asset.Type,
		payout.Asset.Name,
		payout.Asset.Instrument,
		payout.Asset.Shares,
		payout.MonetaryValues.Rate,
		payout.Yield,
		payout.MonetaryValues.Total,
		payout.MonetaryValues.Commission,
		0,
		payout.MonetaryValues.Total,
		0,
		internal.DateTime{Time: payout.BaseTransaction.Timestamp},
	)
}
