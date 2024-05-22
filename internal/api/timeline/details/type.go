//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=type.go -destination type_resolver_mock.go -package=details

package details

import (
	"errors"

	log "github.com/sirupsen/logrus"
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

type TesterFunc func(Response) bool

type TypeResolverInterface interface {
	Resolve(response Response) (Type, error)
}

type TypeResolver struct {
	detectors map[Type]TesterFunc
	logger    *log.Logger
}

func NewTypeResolver(logger *log.Logger) TypeResolver {
	return TypeResolver{
		detectors: map[Type]TesterFunc{
			TypePurchaseTransaction: PurchaseDetector,
			TypeSaleTransaction:     SaleDetector,
			// TypeDepositTransaction:                 DepositDetector,
			// TypeInterestReceivedTransaction: InterestReceivedDetector,
			TypeRoundUpTransaction:        RoundUpDetector,
			TypeSavebackTransaction:       SavebackDetector,
			TypeDividendPayoutTransaction: DividendPayoutDetector,
		},
		logger: logger,
	}
}

func (r TypeResolver) Resolve(response Response) (Type, error) {
	for detectedType, detector := range r.detectors {
		if !detector(response) {
			continue
		}

		r.logger.WithField("id", response.ID).Debugf("%s transaction resolved", detectedType)

		return detectedType, nil
	}

	return TypeUnsupported, ErrUnsupportedResponse
}

func PurchaseDetector(response Response) bool {
	overview, err := response.OverviewSection()
	if err != nil {
		return false
	}

	orderType, err := overview.OrderType()
	if err != nil {
		return false
	}

	return orderType.IsOrderTypePurchase()
}

func SaleDetector(response Response) bool {
	overview, err := response.OverviewSection()
	if err != nil {
		return false
	}

	orderType, err := overview.OrderType()
	if err != nil {
		return false
	}

	return orderType.IsOrderTypeSale()
}

func RoundUpDetector(response Response) bool {
	overview, err := response.OverviewSection()
	if err != nil {
		return false
	}

	orderType, err := overview.OrderType()
	if err != nil {
		return false
	}

	return orderType.IsOrderTypeRoundUp()
}

func SavebackDetector(response Response) bool {
	overview, err := response.OverviewSection()
	if err != nil {
		return false
	}

	orderType, err := overview.OrderType()
	if err != nil {
		return false
	}

	return orderType.IsOrderTypeSaveback()
}

func DepositDetector(response Response) bool {
	overview, err := response.OverviewSection()
	if err != nil {
		return false
	}

	_, err = overview.ReceivedFrom()
	if err == nil {
		return true
	}

	_, err = overview.Deposit()

	return err == nil
}

func InterestReceivedDetector(response Response) bool {
	overview, err := response.OverviewSection()
	if err != nil {
		return false
	}

	_, err = overview.YoY()

	return err == nil
}

func DividendPayoutDetector(response Response) bool {
	overview, err := response.OverviewSection()
	if err != nil {
		return false
	}

	event, err := overview.Event()
	if err != nil {
		return false
	}

	return event.IsEventPayout()
}
