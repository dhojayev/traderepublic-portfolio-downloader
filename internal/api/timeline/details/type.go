//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=type.go -destination type_resolver_mock.go -package=details

package details

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Type int

const (
	TypeUnsupported Type = iota
	TypeSaleTransaction
	TypePurchaseTransaction
	TypeDividendPayoutTransaction
	TypeRoundUpTransaction
	TypeSavebackTransaction
	TypeCardPaymentTransaction
	TypeDepositTransaction
	TypeDepositInterestReceivedTransaction
)

var (
	//nolint: gochecknoglobals
	typeReadableNameMap = map[Type]string{
		TypeUnsupported:                        "unsupported",
		TypeSaleTransaction:                    "sale",
		TypePurchaseTransaction:                "purchase",
		TypeDividendPayoutTransaction:          "dividend payout",
		TypeRoundUpTransaction:                 "round up",
		TypeSavebackTransaction:                "saveback",
		TypeCardPaymentTransaction:             "card payment",
		TypeDepositTransaction:                 "deposit",
		TypeDepositInterestReceivedTransaction: "interest received",
	}

	//nolint: gochecknoglobals
	detectors = map[Type]TesterFunc{
		TypePurchaseTransaction:                PurchaseDetector,
		TypeSaleTransaction:                    SaleDetector,
		// TypeDepositTransaction:                 DepositDetector,
		// TypeDepositInterestReceivedTransaction: InterestReceivedDetector,
		TypeRoundUpTransaction:                 RoundUpDetector,
		TypeSavebackTransaction:                SavebackDetector,
		TypeDividendPayoutTransaction:          DividendPayoutDetector,
	}

	ErrUnsupportedResponse = errors.New("could not resolve transaction type")
)

type TesterFunc func(Response) bool

type TypeResolverInterface interface {
	Resolve(response Response) (Type, error)
}

type TypeResolver struct {
	logger *log.Logger
}

func NewTypeResolver(logger *log.Logger) TypeResolver {
	return TypeResolver{
		logger: logger,
	}
}

func (r TypeResolver) Resolve(response Response) (Type, error) {
	for detectedType, detector := range detectors {
		if !detector(response) {
			continue
		}

		r.logger.WithField("id", response.ID).Debugf("%#v transaction resolved", typeReadableNameMap[detectedType])

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
