//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=type.go -destination type_resolver_mock.go -package=details

package details

import (
	"errors"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
)

type Type string

const (
	TypeUnsupported               Type = "Unsupported"
	TypeSaleTransaction           Type = "Sale"
	TypePurchaseTransaction       Type = "Purchase"
	TypeDividendPayoutTransaction Type = "Dividend payout"
	TypeRoundUpTransaction        Type = "Round up"
	TypeSavebackTransaction       Type = "Saveback"
	TypeCardPaymentTransaction    Type = "Card payment"
	TypeDepositTransaction        Type = "Deposit"
	TypeWithdrawalTransaction     Type = "Withdrawal"
	TypeInterestPayoutTransaction Type = "Interest payout"
)

var ErrTypeResolverUnsupportedType = errors.New("could not resolve transaction type")

type TesterFunc func(transactions.EventType, NormalizedResponse) bool

type TypeResolverInterface interface {
	Resolve(eventType transactions.EventType, response NormalizedResponse) (Type, error)
}

type TypeResolver struct {
	detectors map[Type]TesterFunc
	logger    *log.Logger
}

func NewTypeResolver(logger *log.Logger) TypeResolver {
	return TypeResolver{
		detectors: map[Type]TesterFunc{
			TypeDepositTransaction:        DepositDetector,
			TypeWithdrawalTransaction:     WithdrawalDetector,
			TypeDividendPayoutTransaction: DividendPayoutDetector,
			TypeRoundUpTransaction:        RoundUpDetector,
			TypeSavebackTransaction:       SavebackDetector,
			TypeInterestPayoutTransaction: InterestPayoutDetector,

			// Detectors with the highest performance hit should be listed in the bottom.
			TypePurchaseTransaction: PurchaseDetector,
			TypeSaleTransaction:     SaleDetector,
		},
		logger: logger,
	}
}

func (r TypeResolver) Resolve(eventType transactions.EventType, response NormalizedResponse) (Type, error) {
	for detectedType, detector := range r.detectors {
		if !detector(eventType, response) {
			continue
		}

		r.logger.WithField("id", response.ID).Debugf("%s transaction resolved", detectedType)

		return detectedType, nil
	}

	return Type(""), ErrTypeResolverUnsupportedType
}

func PurchaseDetector(eventType transactions.EventType, response NormalizedResponse) bool {
	supportedEventTypes := []transactions.EventType{
		transactions.EventTypeTradeInvoiceCreated,
		transactions.EventTypeSavingsPlanExecuted,
		transactions.EventTypeSavingsPlanInvoiceCreated,
	}

	if slices.Contains(supportedEventTypes, eventType) {
		return true
	}

	if eventType != transactions.EventTypeOrderExecuted {
		return false
	}

	orderType, err := response.Overview.GetDataByTitles(OverviewDataTitleOrderType)
	if err != nil {
		return false
	}

	return strings.Contains(orderType.Detail.Text, OrderTypeTextsPurchase)
}

func SaleDetector(eventType transactions.EventType, response NormalizedResponse) bool {
	if eventType != transactions.EventTypeOrderExecuted {
		return false
	}

	orderType, err := response.Overview.GetDataByTitles(OverviewDataTitleOrderType)
	if err != nil {
		return false
	}

	return strings.Contains(orderType.Detail.Text, OrderTypeTextsSale)
}

func RoundUpDetector(eventType transactions.EventType, _ NormalizedResponse) bool {
	return eventType == transactions.EventTypeBenefitsSpareChangeExecution
}

func SavebackDetector(eventType transactions.EventType, _ NormalizedResponse) bool {
	return eventType == transactions.EventTypeBenefitsSavebackExecution
}

func DepositDetector(eventType transactions.EventType, _ NormalizedResponse) bool {
	return eventType == transactions.EventTypePaymentInbound ||
		eventType == transactions.EventTypePaymentInboundSepaDirectDebit
}

func InterestPayoutDetector(eventType transactions.EventType, _ NormalizedResponse) bool {
	return eventType == transactions.EventTypeInterestPayoutCreated || eventType == transactions.EventTypeInterestPayout
}

func DividendPayoutDetector(eventType transactions.EventType, _ NormalizedResponse) bool {
	return eventType == transactions.EventTypeCredit || eventType == transactions.EventTypeSSPCorporateActionInvoiceCash
}

func WithdrawalDetector(eventType transactions.EventType, _ NormalizedResponse) bool {
	return eventType == transactions.EventTypePaymentOutbound
}
