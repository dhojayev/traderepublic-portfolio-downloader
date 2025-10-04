package transaction

import (
	"errors"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

var (
	ErrCancelledTransactionReceived = errors.New("canceled transaction type received")
	ErrIgnoredTransactionReceived   = errors.New("ignored transaction type received")
	ErrUnknownTransactionReceived   = errors.New("unknown transaction type received")
)

type Type interface {
	fmt.Stringer
	FindID(details traderepublic.TimelineDetailsJson) string
	FindStatus(details traderepublic.TimelineDetailsJson) (string, error)
	FindTimestamp(details traderepublic.TimelineDetailsJson) (string, error)
	FindISIN(details traderepublic.TimelineDetailsJson) (string, error)
	FindShares(details traderepublic.TimelineDetailsJson) (string, error)
	FindSharePrice(details traderepublic.TimelineDetailsJson) (string, error)
	FindFee(details traderepublic.TimelineDetailsJson) (string, error)
	FindTotal(details traderepublic.TimelineDetailsJson) (string, error)
}

type GenericType struct {
}

func (t *GenericType) FindID(details traderepublic.TimelineDetailsJson) string {
	return string(details.Id)
}

func (t *GenericType) FindStatus(details traderepublic.TimelineDetailsJson) (string, error) {
	header, err := details.SectionHeader()
	if err != nil {
		return "", fmt.Errorf("failed to find header section: %w", err)
	}

	return string(header.Data.Status), nil
}

func (t *GenericType) FindTimestamp(details traderepublic.TimelineDetailsJson) (string, error) {
	header, err := details.SectionHeader()
	if err != nil {
		return "", fmt.Errorf("failed to find header section: %w", err)
	}

	return header.Data.Timestamp, nil
}

type HeaderActionPayloadISINType struct {
	GenericType
}

func (t *HeaderActionPayloadISINType) FindISIN(details traderepublic.TimelineDetailsJson) (string, error) {
	header, err := details.SectionHeader()
	if err != nil {
		return "", fmt.Errorf("failed to find header section: %w", err)
	}

	return header.Action.Payload, nil
}

type SavingsPlanType struct {
	HeaderActionPayloadISINType
}

func (t *SavingsPlanType) FindShares(details traderepublic.TimelineDetailsJson) (string, error) {
	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return "", fmt.Errorf("failed to find overview section: %w", err)
	}

	trn, err := overview.FindData(traderepublic.DataTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to find transaction data: %w", err)
	}

	return *trn.Detail.DisplayValue.Prefix, nil
}

func (t *SavingsPlanType) FindSharePrice(details traderepublic.TimelineDetailsJson) (string, error) {
	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return "", fmt.Errorf("failed to find overview section: %w", err)
	}

	trn, err := overview.FindData(traderepublic.DataTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to find transaction data: %w", err)
	}

	return trn.Detail.DisplayValue.Text, nil
}

func (t *SavingsPlanType) FindFee(details traderepublic.TimelineDetailsJson) (string, error) {
	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return "", fmt.Errorf("failed to find overview section: %w", err)
	}

	fee, err := overview.FindData(traderepublic.DataFee)
	if err != nil {
		return "", fmt.Errorf("failed to find fee data: %w", err)
	}

	return fee.Detail.Text, nil
}

func (t *SavingsPlanType) FindTotal(details traderepublic.TimelineDetailsJson) (string, error) {
	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return "", fmt.Errorf("failed to find overview section: %w", err)
	}

	total, err := overview.FindData(traderepublic.DataTotal)
	if err != nil {
		return "", fmt.Errorf("failed to find total data: %w", err)
	}

	return total.Detail.Text, nil
}

func (t *SavingsPlanType) String() string {
	return "Savings plan"
}

type SavingsPlanPre202502Type struct {
	SavingsPlanType
}

func (t *SavingsPlanPre202502Type) FindShares(details traderepublic.TimelineDetailsJson) (string, error) {
	trnSection, err := details.FindSection(traderepublic.SectionTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to find transaction section: %w", err)
	}

	shares, err := trnSection.FindData(traderepublic.DataShares)
	if err != nil {
		return "", fmt.Errorf("failed to find shares data: %w", err)
	}

	return shares.Detail.Text, nil
}

func (t *SavingsPlanPre202502Type) FindSharePrice(details traderepublic.TimelineDetailsJson) (string, error) {
	trnSection, err := details.FindSection(traderepublic.SectionTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to find transaction section: %w", err)
	}

	sharePrice, err := trnSection.FindData(traderepublic.DataSharePrice)
	if err != nil {
		return "", fmt.Errorf("failed to find shares data: %w", err)
	}

	return sharePrice.Detail.Text, nil
}

func (t *SavingsPlanPre202502Type) FindFee(details traderepublic.TimelineDetailsJson) (string, error) {
	trnSection, err := details.FindSection(traderepublic.SectionTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to find transaction section: %w", err)
	}

	fee, err := trnSection.FindData(traderepublic.DataFee)
	if err != nil {
		return "", fmt.Errorf("failed to find fee data: %w", err)
	}

	return fee.Detail.Text, nil
}

func (t *SavingsPlanPre202502Type) FindTotal(details traderepublic.TimelineDetailsJson) (string, error) {
	trnSection, err := details.FindSection(traderepublic.SectionTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to find transaction section: %w", err)
	}

	total, err := trnSection.FindData(traderepublic.DataTotal)
	if err != nil {
		return "", fmt.Errorf("failed to find total data: %w", err)
	}

	return total.Detail.Text, nil
}

// TransactionType represents the type of a transaction.
type TransactionType string

const (
	TypeUnknown              TransactionType = "unknown"          // Unknown transaction type
	TypeIgnored              TransactionType = "ingored"          // Ignored transaction type
	TypeSavingsplan          TransactionType = "Savings plan"     // Savings plan transaction
	TypeSavingsplanPre202502 TransactionType = "Savings plan"     // Savings plan transaction
	TypeCardPayment          TransactionType = "Card payment"     // Card payment transaction
	TypeCardRefund           TransactionType = "Card refund"      // Card refund transaction
	TypeBuyOrder             TransactionType = "Buy order"        // Buy order transaction
	TypeSellOrder            TransactionType = "Sell order"       // Sell order transaction
	TypeDividendsIncome      TransactionType = "Dividends income" // Dividends income transaction
	TypeRoundUp              TransactionType = "Round up"         // Round up transaction
	TypeSaveback             TransactionType = "Saveback"         // Saveback transaction
	TypeDeposit              TransactionType = "Deposit"          // Deposit transaction
	TypeWithdrawal           TransactionType = "Withdrawal"       // Withdrawal transaction
	TypeInterestPayment      TransactionType = "Interest payment" // Interest payment transaction
)

// TypeResolver resolves the type of a transaction based on its details.
type TypeResolver struct {
}

// NewTypeResolver creates a new instance of TypeResolver.
func NewTypeResolver() *TypeResolver {
	return &TypeResolver{}
}

// Resolve determines the type of a transaction from its details.
func (r *TypeResolver) SetType(details traderepublic.TimelineDetailsJson, model *Model) error {
	header, err := details.SectionHeader()
	if err == nil {
		if header.Data.Status == "canceled" {
			return ErrCancelledTransactionReceived
		}
	}

	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return fmt.Errorf("failed to find overview section: %w", err)
	}

	// Check for ignored transactions
	_, err = overview.FindData(traderepublic.DataCardVerification)
	if err == nil {
		return fmt.Errorf("%w: %s", ErrIgnoredTransactionReceived, details.Id)
	}

	// // Check for card payment transactions
	// _, err = overview.FindData(traderepublic.DataCardPayment)
	// if err == nil {
	// 	model.Type = TypeCardPayment

	// 	return
	// }

	// // Check for card refund transactions
	// _, err = overview.FindData(traderepublic.DataCardRefund)
	// if err == nil {
	// 	model.Type = TypeCardRefund

	// 	return
	// }

	// // Check for deposit transactions
	// _, err = details.FindSection(traderepublic.SectionSender)
	// if err == nil {
	// 	model.Type = TypeDeposit

	// 	return
	// }
	// _, err = overview.FindData(traderepublic.DataFrom)
	// if err == nil {
	// 	model.Type = TypeDeposit

	// 	return
	// }

	// // Check for withdrawal transactions
	// _, err = overview.FindData(traderepublic.DataTo)
	// if err == nil {
	// 	model.Type = TypeWithdrawal

	// 	return
	// }

	// payment, err := overview.FindData(traderepublic.DataPayment)
	// if err == nil {
	// 	if payment.Detail.Text == "Direct Debit" {
	// 		model.Type = TypeDeposit

	// 		return
	// 	}
	// }

	// Check for savings plan transactions
	_, err = overview.FindData(traderepublic.DataSavingsPlan)
	if err == nil {
		model.Type = &SavingsPlanType{}

		return nil
	}

	// // Check for dividends income transactions
	// event, err := overview.FindData(traderepublic.DataEvent)
	// if err == nil {
	// 	switch event.Detail.Text {
	// 	case "Income", "Cash dividend":
	// 		model.Type = TypeDividendsIncome

	// 		return
	// 	case "Tax Settlement":
	// 		model.Type = TypeIgnored

	// 		slog.Warn("ignored transaction type", "id", details.Id, "err", err, "info", "Tax Settlement")

	// 		return
	// 	}
	// }

	// Check for buy order transactions
	orderType, err := overview.FindData(traderepublic.DataOrderType)
	if err == nil {
		switch orderType.Detail.Text {
		case "Savings plan":
			model.Type = &SavingsPlanPre202502Type{}

			return nil
			// 	case "Buy":
			// 		model.Type = TypeBuyOrder

			// 		return
			// 	case "Sell":
			// 		model.Type = TypeSellOrder

			// 		return
			// 	case "Limit Sell":
			// 		model.Type = TypeSellOrder

			// 		return
		}
	}

	// // Check for interest payment transactions
	// _, err = overview.FindData(traderepublic.DataAverageBalance)
	// if err == nil {
	// 	model.Type = TypeInterestPayment

	// 	return
	// }
	// steps, err := details.SectionSteps()
	// if err == nil {
	// 	_, err = steps.FindStep(traderepublic.StepInterestPayment)
	// 	if err == nil {
	// 		model.Type = TypeInterestPayment

	// 		return
	// 	}
	// }

	// // Check for saveback transactions
	// _, err = overview.FindData(traderepublic.DataSaveback)
	// if err == nil {
	// 	model.Type = TypeSaveback

	// 	return
	// }

	// // Check for round up transactions
	// _, err = overview.FindData(traderepublic.DataRoundUp)
	// if err == nil {
	// 	model.Type = TypeRoundUp

	// 	return
	// }

	// // Check for limit sell transactions
	// _, err = overview.FindData(traderepublic.DataLimitSell)
	// if err == nil {
	// 	model.Type = TypeSellOrder

	// 	return
	// }

	// // Check for sell transactions
	// _, err = overview.FindData(traderepublic.DataSell)
	// if err == nil {
	// 	model.Type = TypeSellOrder

	// 	return
	// }

	// // Check for buy transactions
	// _, err = overview.FindData(traderepublic.DataBuy)
	// if err == nil {
	// 	model.Type = TypeBuyOrder

	// 	return
	// }

	return fmt.Errorf("%w: %s", ErrUnknownTransactionReceived, details.Id)
}
