package transaction

import (
	"errors"
	"fmt"
	"slices"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	log "github.com/sirupsen/logrus"
)

type HandlerInterface interface {
	Handle() error
}

type Handler struct {
	listClient        transactions.ClientInterface
	detailsClient     details.ClientInterface
	eventTypeResolver transactions.EventTypeResolverInterface
	processor         ProcessorInterface
	logger            *log.Logger
}

func NewHandler(
	listClient transactions.ClientInterface,
	detailsClient details.ClientInterface,
	eventTypeResolver transactions.EventTypeResolverInterface,
	processor ProcessorInterface,
	logger *log.Logger,
) Handler {
	return Handler{
		listClient:        listClient,
		detailsClient:     detailsClient,
		eventTypeResolver: eventTypeResolver,
		processor:         processor,
		logger:            logger,
	}
}

func (h Handler) Handle() error {
	counter := internal.NewOperationCounter()

	responses, err := h.GetTimelineTransactions()
	if err != nil {
		return err
	}

	for _, response := range responses {
		if !response.Action.HasDetails() {
			continue
		}

		infoFields := log.Fields{"id": response.ID}

		err := h.ProcessTransactionResponse(response)

		// Handle ignorable errors.
		switch {
		case err == nil:
			counter.Processed().Add(1)
		case errors.Is(err, websocket.ErrMsgErrorStateReceived):
			h.logger.WithFields(infoFields).Error(err)
			counter.Skipped().Add(1)

			continue
		case errors.Is(err, transactions.ErrEventTypeUnsupported):
			h.logger.WithFields(infoFields).Warn("Unsupported transaction skipped")
			counter.Skipped().Add(1)

			continue
		case errors.Is(err, ErrModelBuilderUnsupportedType):
			h.logger.WithFields(infoFields).Warn("Unsupported transaction skipped")
			counter.Skipped().Add(1)

			continue
		case errors.Is(err, ErrModelBuilderInsufficientDataResolved):
			h.logger.WithFields(infoFields).Warnf("Transaction skipped due to missing details: %s", err)
			counter.Skipped().Add(1)

			continue
		}
	}

	h.logger.Infof(
		"Transactions completed: %d, skipped: %d",
		counter.Processed().Load(),
		counter.Skipped().Load(),
	)

	return nil
}

func (h Handler) GetTimelineTransactions() ([]transactions.ResponseItem, error) {
	h.logger.Info("Downloading transactions")

	var transactions []transactions.ResponseItem

	err := h.listClient.List(&transactions)
	if err != nil {
		return transactions, fmt.Errorf("could not fetch transactions: %w", err)
	}

	slices.Reverse(transactions)

	h.logger.Infof("%d transactions downloaded", len(transactions))

	return transactions, nil
}

func (h Handler) ProcessTransactionResponse(transaction transactions.ResponseItem) error {
	infoFields := log.Fields{"id": transaction.ID}

	h.logger.WithFields(infoFields).Info("Fetching transaction details")

	var details details.Response

	err := h.detailsClient.Details(transaction.Action.Payload, &details)
	if err != nil {
		return fmt.Errorf("could not fetch transaction details: %w", err)
	}

	eventType, err := h.eventTypeResolver.Resolve(transaction)
	if err != nil {
		return fmt.Errorf("could not resolve transaction even type: %w", err)
	}

	h.logger.WithFields(infoFields).Info("Processing transaction details")

	if err := h.processor.Process(eventType, details); err != nil {
		return fmt.Errorf("could not process transaction: %w", err)
	}

	h.logger.WithFields(infoFields).Info("Transaction processed")

	return nil
}
