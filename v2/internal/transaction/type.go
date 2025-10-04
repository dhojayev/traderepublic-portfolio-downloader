package transaction

import (
	"fmt"
	"slices"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type TransactionType string

const (
	TypeUnknown         TransactionType = "unknown"
	TypeIgnored         TransactionType = "ingored"
	TypeSavingsplan     TransactionType = "Savings plan"
	TypeCardPayment     TransactionType = "Card payment"
	TypeCardRefund      TransactionType = "Card refund"
	TypeBuyOrder        TransactionType = "Buy order"
	TypeSellOrder       TransactionType = "Sell order"
	TypeDividendsIncome TransactionType = "Dividends income"
	TypeRoundUp         TransactionType = "Round up"
	TypeSaveback        TransactionType = "Saveback"
	TypeDeposit         TransactionType = "Deposit"
	TypeWithdrawal      TransactionType = "Withdrawal"
	TypeInterestPayment TransactionType = "Interest payment"
)

type TypeResolver struct {
}

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{}
}

func (r *TypeResolver) Resolve(details traderepublic.TimelineDetailsJson) (TransactionType, error) {
	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return TypeUnknown, fmt.Errorf("unknown transaction type id %s: %w", details.Id, err)
	}

	_, err = overview.FindData(traderepublic.DataCardVerification)
	if err == nil {
		return TypeIgnored, nil
	}

	_, err = overview.FindData(traderepublic.DataCardPayment)
	if err == nil {
		return TypeCardPayment, nil
	}

	_, err = overview.FindData(traderepublic.DataCardRefund)
	if err == nil {
		return TypeCardRefund, nil
	}

	_, err = details.FindSection(traderepublic.SectionSender)
	if err == nil {
		return TypeDeposit, nil
	}

	_, err = overview.FindData(traderepublic.DataFrom)
	if err == nil {
		return TypeDeposit, nil
	}

	_, err = overview.FindData(traderepublic.DataTo)
	if err == nil {
		return TypeWithdrawal, nil
	}

	_, err = overview.FindData(traderepublic.DataPayment)
	if err == nil {
		return TypeDeposit, nil
	}

	_, err = overview.FindData(traderepublic.DataSavingsPlan)
	if err == nil {
		return TypeSavingsplan, nil
	}

	event, err := overview.FindData(traderepublic.DataEvent)
	if err == nil {
		texts := []string{"Income", "Cash dividend"}
		if slices.Contains(texts, event.Detail.Text) {
			return TypeDividendsIncome, nil
		}
	}

	orderType, err := overview.FindData(traderepublic.DataOrderType)
	if err == nil {
		switch orderType.Detail.Text {
		case "Savings plan":
			return TypeSavingsplan, nil
		case "Buy":
			return TypeBuyOrder, nil
		case "Sell":
			return TypeSellOrder, nil
		case "Limit Sell":
			return TypeSellOrder, nil
		}
	}

	_, err = overview.FindData(traderepublic.DataAverageBalance)
	if err == nil {
		return TypeInterestPayment, nil
	}

	steps, err := details.SectionSteps()
	if err == nil {
		_, err = steps.FindStep(traderepublic.StepInterestPayment)
		if err == nil {
			return TypeInterestPayment, nil
		}
	}

	_, err = overview.FindData(traderepublic.DataSaveback)
	if err == nil {
		return TypeSaveback, nil
	}

	_, err = overview.FindData(traderepublic.DataRoundUp)
	if err == nil {
		return TypeRoundUp, nil
	}

	_, err = overview.FindData(traderepublic.DataLimitSell)
	if err == nil {
		return TypeSellOrder, nil
	}

	_, err = overview.FindData(traderepublic.DataSell)
	if err == nil {
		return TypeSellOrder, nil
	}

	_, err = overview.FindData(traderepublic.DataBuy)
	if err == nil {
		return TypeBuyOrder, nil
	}

	return TypeUnknown, fmt.Errorf("unknown transaction type id: %s", details.Id)
}
