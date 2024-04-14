package transaction

import (
	"errors"
	"fmt"

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

func (f CSVEntryFactory) Make(transaction Model) (filesystem.CSVEntry, error) {
	var debit, credit, investedAmount float64

	yield := transaction.Yield
	profit := transaction.Profit
	shares := transaction.Shares
	rate := transaction.Rate
	commission := transaction.Commission

	switch transaction.Type {
	case TypePurchase:
		debit = transaction.Total
		investedAmount = transaction.Total - transaction.Commission
	case TypeSale:
		shares = -shares
		credit = transaction.Total
		investedAmount = -(transaction.Total - transaction.Profit + transaction.Commission)
	case TypeSaveback:

	case TypeRoundUp:
		debit = transaction.Total
	case TypeDividendPayout:
		profit = transaction.Total
		credit = transaction.Total
	default:
		return filesystem.CSVEntry{}, fmt.Errorf(
			"unsupported type '%s' received: %w",
			transaction.Type,
			errors.New("unsupported value object received"),
		)
	}

	return filesystem.NewCSVEntry(
		transaction.UUID,
		transaction.Status,
		transaction.Type,
		transaction.Instrument.Type(),
		transaction.Instrument.Name,
		transaction.Instrument.ISIN,
		shares,
		rate,
		yield,
		profit,
		commission,
		debit,
		credit,
		investedAmount,
		internal.DateTime{Time: transaction.Timestamp},
	), nil
}
