//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=type.go -destination type_resolver_mock.go -package=details

package details

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
)

type Type string

const (
	TypeUnsupported                 Type = "Unsupported"
	TypeSaleTransaction             Type = "Sale"
	TypePurchaseTransaction         Type = "Purchase"
	TypeDividendPayoutTransaction   Type = "Dividend payout"
	TypeRoundUpTransaction          Type = "Round up"
	TypeSavebackTransaction         Type = "Saveback"
	TypeCardPaymentTransaction      Type = "Card payment"
	TypeDepositTransaction          Type = "Deposit"
	TypeInterestReceivedTransaction Type = "Interest received"
)

var ErrUnsupportedResponse = errors.New("could not resolve transaction type")

type TesterFunc func(transactions.EventType, Response) bool

type TypeResolverInterface interface {
	Resolve(eventType transactions.EventType, response Response) (Type, error)
}

type TypeResolver struct {
	detectors map[Type]TesterFunc
	logger    *log.Logger
}

func NewTypeResolver(logger *log.Logger) TypeResolver {
	return TypeResolver{
		detectors: map[Type]TesterFunc{
			TypeDepositTransaction:          DepositDetector,
			TypeInterestReceivedTransaction: InterestReceivedDetector,
			TypeRoundUpTransaction:          RoundUpDetector,
			TypeSavebackTransaction:         SavebackDetector,
			TypeDividendPayoutTransaction:   DividendPayoutDetector,

			// Detectors with the highest performance hit should be listed in the bottom.
			TypePurchaseTransaction: PurchaseDetector,
			TypeSaleTransaction:     SaleDetector,
		},
		logger: logger,
	}
}

func (r TypeResolver) Resolve(eventType transactions.EventType, response Response) (Type, error) {
	for detectedType, detector := range r.detectors {
		if !detector(eventType, response) {
			continue
		}

		r.logger.WithField("id", response.ID).Debugf("%s transaction resolved", detectedType)

		return detectedType, nil
	}

	return TypeUnsupported, ErrUnsupportedResponse
}

func PurchaseDetector(eventType transactions.EventType, response Response) bool {
	if eventType == transactions.EvenTypeSavingsPlanExecuted {
		return true
	}

	if eventType != transactions.EventTypeOrderExecuted {
		return false
	}

	tableSections, err := response.SectionsTypeTable()
	if err != nil {
		return false
	}

	overviewSection, err := tableSections.FindByTitle(SectionTitleOverview)
	if err != nil {
		return false
	}

	orderType, err := overviewSection.GetDataByTitle(overviewDataTitleOrderType)
	if err != nil {
		return false
	}

	return strings.Contains(orderType.Detail.Text, orderTypeTextsPurchase)
}

func SaleDetector(eventType transactions.EventType, response Response) bool {
	if eventType != transactions.EventTypeOrderExecuted {
		return false
	}

	tableSections, err := response.SectionsTypeTable()
	if err != nil {
		return false
	}

	overviewSection, err := tableSections.FindByTitle(SectionTitleOverview)
	if err != nil {
		return false
	}

	orderType, err := overviewSection.GetDataByTitle(overviewDataTitleOrderType)
	if err != nil {
		return false
	}

	return strings.Contains(orderType.Detail.Text, orderTypeTextsSale)
}

func RoundUpDetector(eventType transactions.EventType, _ Response) bool {
	return eventType == transactions.EventTypeBenefitsSpareChangeExecution
}

func SavebackDetector(eventType transactions.EventType, _ Response) bool {
	return eventType == transactions.EventTypeBenefitsSavebackExecution
}

func DepositDetector(eventType transactions.EventType, _ Response) bool {
	return eventType == transactions.EventTypePaymentInbound ||
		eventType == transactions.EventTypePaymentInboundSepaDirectDebit
}

func InterestReceivedDetector(eventType transactions.EventType, _ Response) bool {
	return eventType == transactions.EventTypeInterestPayoutCreated
}

func DividendPayoutDetector(eventType transactions.EventType, _ Response) bool {
	return eventType == transactions.EvenTypeCredit
}
