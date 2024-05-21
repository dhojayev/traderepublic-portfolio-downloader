//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=resolver.go -destination resolver_mock.go -package=transaction

package transaction

import (
	"errors"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"

	log "github.com/sirupsen/logrus"
)

var ErrUnsupportedResponse = errors.New("could not resolve transaction type")

type TypeResolverInterface interface {
	Resolve(response details.Response) (details.Type, error)
}

type TypeResolver struct {
	logger *log.Logger
}

func NewTypeResolver(logger *log.Logger) TypeResolver {
	return TypeResolver{
		logger: logger,
	}
}

func (r TypeResolver) Resolve(response details.Response) (details.Type, error) {
	overview, err := response.OverviewSection()
	if err != nil {
		return details.TypeUnsupported, fmt.Errorf("response error: %w", err)
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

		return details.TypeDepositTransaction, nil
	case yoyErr == nil:
		r.logger.WithFields(logFields).Debug("interest received transaction resolved")

		return details.TypeDepositInterestReceivedTransaction, nil
	case orderTypeErr == nil:
		return r.ResolveByOrderType(orderType, logFields)
	case eventErr == nil:
		if !event.IsEventPayout() {
			return details.TypeUnsupported, fmt.Errorf("%w: %w", ErrUnsupportedResponse, err)
		}

		r.logger.WithFields(logFields).Debug("purchase transaction resolved")

		return details.TypeDividendPayoutTransaction, nil
	}

	return details.TypeUnsupported, ErrUnsupportedResponse
}

func (r TypeResolver) ResolveByOrderType(
	order details.ResponseSectionTypeTableData,
	logFields log.Fields,
) (details.Type, error) {
	switch {
	case order.IsOrderTypeSale():
		r.logger.WithFields(logFields).Debug("sale transaction resolved")

		return details.TypeSaleTransaction, nil
	case order.IsOrderTypePurchase():
		r.logger.WithFields(logFields).Debug("purchase transaction resolved")

		return details.TypePurchaseTransaction, nil
	case order.IsOrderTypeRoundUp():
		r.logger.WithFields(logFields).Debug("round up transaction resolved")

		return details.TypeRoundUpTransaction, nil
	case order.IsOrderTypeSaveback():
		r.logger.WithFields(logFields).Debug("saveback transaction resolved")

		return details.TypeSavebackTransaction, nil
	}

	return details.TypeUnsupported, errors.New("could not resolve transaction type from order")
}
