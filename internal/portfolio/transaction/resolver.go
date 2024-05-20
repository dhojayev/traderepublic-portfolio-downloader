//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=resolver.go -destination resolver_mock.go -package=transaction

package transaction

import (
	"errors"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"

	log "github.com/sirupsen/logrus"
)

var ErrUnsupportedResponse = errors.New("could not resolve transaction type")

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

type Type int

type TypeResolverInterface interface {
	Resolve(response details.Response) (Type, error)
}

type TypeResolver struct {
	logger *log.Logger
}

func NewTypeResolver(logger *log.Logger) TypeResolver {
	return TypeResolver{
		logger: logger,
	}
}

func (r TypeResolver) Resolve(response details.Response) (Type, error) {
	overview, err := response.OverviewSection()
	if err != nil {
		return TypeUnsupported, fmt.Errorf("response error: %w", err)
	}

	logFields := log.Fields{"id": response.ID}

	_, receivedFromErr := overview.ReceivedFrom()
	_, depositErr := overview.Deposit()
	_, yoyErr := overview.YoY()
	orderType, orderTypeErr := overview.OrderType()
	event, eventErr := overview.Event()

	switch {
	case receivedFromErr == nil, depositErr == nil:
		r.logger.WithFields(logFields).Debug("deposit transaction resolved")

		return TypeDepositTransaction, nil
	case yoyErr == nil:
		r.logger.WithFields(logFields).Debug("interest received transaction resolved")

		return TypeDepositInterestReceivedTransaction, nil
	case orderTypeErr == nil:
		return r.ResolveByOrderType(orderType, logFields)
	case eventErr == nil:
		if !event.IsEventPayout() {
			return TypeUnsupported, fmt.Errorf("%w: %w", ErrUnsupportedResponse, err)
		}

		r.logger.WithFields(logFields).Debug("purchase transaction resolved")

		return TypeDividendPayoutTransaction, nil
	}

	return TypeUnsupported, ErrUnsupportedResponse
}

func (r TypeResolver) ResolveByOrderType(order details.ResponseSectionTypeTableData, logFields log.Fields) (Type, error) {
	switch {
	case order.IsOrderTypeSale():
		r.logger.WithFields(logFields).Debug("sale transaction resolved")

		return TypeSaleTransaction, nil
	case order.IsOrderTypePurchase():
		r.logger.WithFields(logFields).Debug("purchase transaction resolved")

		return TypePurchaseTransaction, nil
	case order.IsOrderTypeRoundUp():
		r.logger.WithFields(logFields).Debug("round up transaction resolved")

		return TypeRoundUpTransaction, nil
	case order.IsOrderTypeSaveback():
		r.logger.WithFields(logFields).Debug("saveback transaction resolved")

		return TypeSavebackTransaction, nil
	}

	return TypeUnsupported, errors.New("could not resolve transaction type from order")
}
