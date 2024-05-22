package transactions

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

const (
	EventTypePaymentInbound                EventType = "PAYMENT_INBOUND"
	EventTypePaymentInboundSepaDirectDebit EventType = "PAYMENT_INBOUND_SEPA_DIRECT_DEBIT"
	EventTypeOrderExecuted                 EventType = "ORDER_EXECUTED"
	EvenTypeSavingsPlanExecuted            EventType = "SAVINGS_PLAN_EXECUTED"
	EventTypeInterestPayoutCreated         EventType = "INTEREST_PAYOUT_CREATED"
	EvenTypeCredit                         EventType = "CREDIT"
	EventTypeBenefitsSavebackExecution     EventType = "benefits_saveback_execution"
	EventTypeBenefitsSpareChangeExecution  EventType = "benefits_spare_change_execution"
	EventTypeCardSuccessfulTransaction     EventType = "card_successful_transaction"
	EventTypeCardRefund                    EventType = "card_refund"
)

var ErrUnsupportedEventType = errors.New("unsupported event type")

type EventType string

type EventTypeResolverInterface interface {
	Resolve(response TransactionResponse) (EventType, error)
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
			EventTypeOrderExecuted,
			EvenTypeSavingsPlanExecuted,
			EventTypeInterestPayoutCreated,
			EvenTypeCredit,
			EventTypeBenefitsSavebackExecution,
			EventTypeBenefitsSpareChangeExecution,
		},
		logger: logger,
	}
}

func (e EventTypeResolver) Resolve(response TransactionResponse) (EventType, error) {
	for _, t := range e.supportedTypes {
		if response.EventType != string(t) {
			continue
		}

		return t, nil
	}

	return "", ErrUnsupportedEventType
}
