package transaction

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

// TransactionType represents the type of a transaction.
type TransactionType string

const (
	TypeUnknown         TransactionType = "unknown"          // Unknown transaction type
	TypeIgnored         TransactionType = "ingored"          // Ignored transaction type
	TypeSavingsplan     TransactionType = "Savings plan"     // Savings plan transaction
	TypeCardPayment     TransactionType = "Card payment"     // Card payment transaction
	TypeCardRefund      TransactionType = "Card refund"      // Card refund transaction
	TypeBuyOrder        TransactionType = "Buy order"        // Buy order transaction
	TypeSellOrder       TransactionType = "Sell order"       // Sell order transaction
	TypeDividendsIncome TransactionType = "Dividends income" // Dividends income transaction
	TypeRoundUp         TransactionType = "Round up"         // Round up transaction
	TypeSaveback        TransactionType = "Saveback"         // Saveback transaction
	TypeDeposit         TransactionType = "Deposit"          // Deposit transaction
	TypeWithdrawal      TransactionType = "Withdrawal"       // Withdrawal transaction
	TypeInterestPayment TransactionType = "Interest payment" // Interest payment transaction
)

// TypeResolver resolves the type of a transaction based on its details.
type TypeResolver struct {
}

// NewTypeResolver creates a new instance of TypeResolver.
func NewTypeResolver() *TypeResolver {
	return &TypeResolver{}
}

// Resolve determines the type of a transaction from its details.
func (r *TypeResolver) Resolve(details traderepublic.TimelineDetailsJson) (TransactionType, error) {
	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return TypeUnknown, fmt.Errorf("unknown transaction type id %s: %w", details.Id, err)
	}

	// Check for ignored transactions
	_, err = overview.FindData(traderepublic.DataCardVerification)
	if err == nil {
		return TypeIgnored, nil
	}

	// Check for card payment transactions
	_, err = overview.FindData(traderepublic.DataCardPayment)
	if err == nil {
		return TypeCardPayment, nil
	}

	// Check for card refund transactions
	_, err = overview.FindData(traderepublic.DataCardRefund)
	if err == nil {
		return TypeCardRefund, nil
	}

	// Check for deposit transactions
	_, err = details.FindSection(traderepublic.SectionSender)
	if err == nil {
		return TypeDeposit, nil
	}
	_, err = overview.FindData(traderepublic.DataFrom)
	if err == nil {
		return TypeDeposit, nil
	}

	// Check for withdrawal transactions
	_, err = overview.FindData(traderepublic.DataTo)
	if err == nil {
		return TypeWithdrawal, nil
	}

	payment, err := overview.FindData(traderepublic.DataPayment)
	if err == nil {
		if payment.Detail.Text == "Direct Debit" {
			return TypeDeposit, nil
		}
	}

	// Check for savings plan transactions
	_, err = overview.FindData(traderepublic.DataSavingsPlan)
	if err == nil {
		return TypeSavingsplan, nil
	}

	// Check for dividends income transactions
	event, err := overview.FindData(traderepublic.DataEvent)
	if err == nil {
		switch event.Detail.Text {
		case "Income", "Cash dividend":
			return TypeDividendsIncome, nil
		case "Tax Settlement":
			return TypeIgnored, nil
		}
	}

	// Check for buy order transactions
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

	// Check for interest payment transactions
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

	// Check for saveback transactions
	_, err = overview.FindData(traderepublic.DataSaveback)
	if err == nil {
		return TypeSaveback, nil
	}

	// Check for round up transactions
	_, err = overview.FindData(traderepublic.DataRoundUp)
	if err == nil {
		return TypeRoundUp, nil
	}

	// Check for limit sell transactions
	_, err = overview.FindData(traderepublic.DataLimitSell)
	if err == nil {
		return TypeSellOrder, nil
	}

	// Check for sell transactions
	_, err = overview.FindData(traderepublic.DataSell)
	if err == nil {
		return TypeSellOrder, nil
	}

	// Check for buy transactions
	_, err = overview.FindData(traderepublic.DataBuy)
	if err == nil {
		return TypeBuyOrder, nil
	}

	return TypeUnknown, fmt.Errorf("unknown transaction type id: %s", details.Id)
}
