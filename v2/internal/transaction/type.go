package transaction

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type TransactionType string

const (
	TypeUnknown        TransactionType = "unknown"
	TypeSavingsplan    TransactionType = "Savings plan"
	TypePurchase       TransactionType = "Purchase"
	TypeSale           TransactionType = "Sale"
	TypeDividendPayout TransactionType = "Dividends"
	TypeRoundUp        TransactionType = "Round up"
	TypeSaveback       TransactionType = "Saveback"
	TypeDeposit        TransactionType = "Deposit"
	TypeWithdrawal     TransactionType = "Withdrawal"
	TypeInterestPayout TransactionType = "Interest payout"
)

type TypeResolver struct {
}

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{}
}

func (r *TypeResolver) Resolve(details traderepublic.TimelineDetailsJson) (TransactionType, error) {
	switch {
	case r.SavingsPlan(details):
		return TypeSavingsplan, nil
	}

	return TypeUnknown, fmt.Errorf("unknown transaction type id: %s", details.Id)
}

func (r *TypeResolver) SavingsPlan(details traderepublic.TimelineDetailsJson) bool {
	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return false
	}

	orderType, err := overview.FindData(traderepublic.DataOrderType)
	if err != nil {
		_, err := overview.FindData(traderepublic.DataSavingsPlan)
		if err != nil {
			return false
		}
	}

	return orderType.Detail.Text == "Savings plan"
}
