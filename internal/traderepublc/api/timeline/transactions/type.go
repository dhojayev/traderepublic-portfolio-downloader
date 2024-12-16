//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=type.go -destination type_resolver_mock.go -package=transactions

package transactions

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

const (
	EventTypePaymentInbound                EventType = "PAYMENT_INBOUND"
	EventTypePaymentInboundSepaDirectDebit EventType = "PAYMENT_INBOUND_SEPA_DIRECT_DEBIT"
	EventTypePaymentOutbound               EventType = "PAYMENT_OUTBOUND"
	EventTypeOrderExecuted                 EventType = "ORDER_EXECUTED"
	EventTypeSavingsPlanExecuted           EventType = "SAVINGS_PLAN_EXECUTED"
	EventTypeSavingsPlanInvoiceCreated     EventType = "SAVINGS_PLAN_INVOICE_CREATED"
	EventTypeInterestPayoutCreated         EventType = "INTEREST_PAYOUT_CREATED"
	EventTypeInterestPayout                EventType = "INTEREST_PAYOUT"
	EventTypeCredit                        EventType = "CREDIT"
	EventTypeBenefitsSavebackExecution     EventType = "benefits_saveback_execution"
	EventTypeBenefitsSpareChangeExecution  EventType = "benefits_spare_change_execution"
	EventTypeSSPCorporateActionInvoiceCash EventType = "ssp_corporate_action_invoice_cash"
	EventTypeCardSuccessfulTransaction     EventType = "card_successful_transaction"
	EventTypeCardRefund                    EventType = "card_refund"
)

var ErrEventTypeUnsupported = errors.New("unsupported event type")

type EventType string

type EventTypeResolverInterface interface {
	Resolve(response ResponseItem) (EventType, error)
}

type EventTypeResolver struct {
	supportedTypes []EventType
	logger         *log.Logger
}

func NewEventTypeResolver(logger *log.Logger) EventTypeResolver {
	return EventTypeResolver{
		supportedTypes: []EventType{
			EventTypePaymentInbound,
			EventTypePaymentInboundSepaDirectDebit,
			EventTypePaymentOutbound,
			EventTypeOrderExecuted,
			EventTypeSavingsPlanExecuted,
			EventTypeSavingsPlanInvoiceCreated,
			EventTypeInterestPayoutCreated,
			EventTypeInterestPayout,
			EventTypeCredit,
			EventTypeBenefitsSavebackExecution,
			EventTypeBenefitsSpareChangeExecution,
			EventTypeSSPCorporateActionInvoiceCash,
		},
		logger: logger,
	}
}

func (e EventTypeResolver) Resolve(response ResponseItem) (EventType, error) {
	for _, t := range e.supportedTypes {
		if response.EventType != string(t) {
			continue
		}

		return t, nil
	}

	return "", ErrEventTypeUnsupported
}
