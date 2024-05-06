//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=resolver.go -destination resolver_mock.go -package=transaction

package transaction

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"

	log "github.com/sirupsen/logrus"
)

const (
	TypeUnsupported Type = iota
	TypeSaleTransaction
	TypePurchaseTransaction
	TypeDividendPayoutTransaction
	TypeRoundUpTransaction
	TypeSavebackTransaction
	TypeCardPaymentTransaction
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

	orderType, err := overview.OrderType()
	if err != nil {
		event, err := overview.Event()
		if err != nil {
			return TypeUnsupported, fmt.Errorf("overview contains no order type nor event: %w: %w", ErrUnsupportedResponse, err)
		}

		if !event.IsEventPayout() {
			return TypeUnsupported, fmt.Errorf("%w: %w", ErrUnsupportedResponse, err)
		}

		r.logger.WithFields(logFields).Debug("purchase transaction resolved")

		return TypeDividendPayoutTransaction, nil
	}

	switch {
	case orderType.IsOrderTypeSale():
		r.logger.WithFields(logFields).Debug("sale transaction resolved")

		return TypeSaleTransaction, nil
	case orderType.IsOrderTypePurchase():
		r.logger.WithFields(logFields).Debug("purchase transaction resolved")

		return TypePurchaseTransaction, nil
	case orderType.IsOrderTypeRoundUp():
		r.logger.WithFields(logFields).Debug("round up transaction resolved")

		return TypeRoundUpTransaction, nil
	case orderType.IsOrderTypeSaveback():
		r.logger.WithFields(logFields).Debug("saveback transaction resolved")

		return TypeSavebackTransaction, nil
	}

	return TypeUnsupported, fmt.Errorf("could not resolve transaction type: %w", err)
}
