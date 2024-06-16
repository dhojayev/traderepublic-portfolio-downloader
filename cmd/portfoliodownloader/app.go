package portfoliodownloader

import (
	"errors"
	"fmt"
	"slices"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/websocket"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

type App struct {
	transactionsClient    transactions.ClientInterface
	eventTypeResolver     transactions.EventTypeResolverInterface
	timelineDetailsClient details.ClientInterface
	transactionProcessor  transaction.ProcessorInterface
	activityLogClient     activitylog.ClientInterface
	logger                *log.Logger
}

func NewApp(
	transactionsClient transactions.ClientInterface,
	eventTypeResolver transactions.EventTypeResolverInterface,
	timelineDetailsClient details.ClientInterface,
	transactionProcessor transaction.ProcessorInterface,
	activityLogClient activitylog.ClientInterface,
	logger *log.Logger,
) App {
	return App{
		transactionsClient:    transactionsClient,
		eventTypeResolver:     eventTypeResolver,
		timelineDetailsClient: timelineDetailsClient,
		transactionProcessor:  transactionProcessor,
		activityLogClient:     activityLogClient,
		logger:                logger,
	}
}

//nolint:funlen,cyclop
func (a App) Run() error {
	counter := 0

	entries, err := a.GetActivityLog()
	if err != nil {
		return fmt.Errorf("could not fetch activity log entries: %w", err)
	}

	for _, entry := range entries {
		if !entry.Action.HasDetails() {
			continue
		}

		infoFields := log.Fields{"id": entry.ID}

		a.logger.WithFields(infoFields).Info("Fetching activity entry details")

		var details details.Response

		err := a.timelineDetailsClient.Details(entry.Action.Payload, &details)
		if err != nil {
			return fmt.Errorf("could not fetch activity entry details: %w", err)
		}
	}

	responses, err := a.GetTimelineTransactions()
	if err != nil {
		return err
	}

	for _, response := range responses {
		if !response.Action.HasDetails() {
			continue
		}

		infoFields := log.Fields{"id": response.ID}

		err := a.ProcessTransactionResponse(response)

		// Handle ignorable errors.
		switch {
		case err == nil:
			counter++
		case errors.Is(err, websocket.ErrMsgErrorStateReceived):
			a.logger.WithFields(infoFields).Error(err)

			continue
		case errors.Is(err, transactions.ErrEventTypeUnsupported):
			a.logger.WithFields(infoFields).Warn("Unsupported transaction skipped")

			continue
		case errors.Is(err, transaction.ErrModelBuilderUnsupportedType):
			a.logger.WithFields(infoFields).Warn("Unsupported transaction skipped")

			continue
		case errors.Is(err, transaction.ErrModelBuilderInsufficientDataResolved):
			a.logger.WithFields(infoFields).Warnf("Transaction skipped due to missing details: %s", err)

			continue
		}
	}

	skippedCount := len(responses) - counter

	a.logger.Infof("Completed: %d, skipped: %d", counter, skippedCount)

	return nil
}

func (a App) GetTimelineTransactions() ([]transactions.ResponseItem, error) {
	a.logger.Info("Downloading transactions")

	var transactions []transactions.ResponseItem

	err := a.transactionsClient.List(&transactions)
	if err != nil {
		return transactions, fmt.Errorf("could not fetch transactions: %w", err)
	}

	slices.Reverse(transactions)

	a.logger.Infof("%d transactions downloaded", len(transactions))

	return transactions, nil
}

func (a App) ProcessTransactionResponse(transaction transactions.ResponseItem) error {
	infoFields := log.Fields{"id": transaction.ID}

	a.logger.WithFields(infoFields).Info("Fetching transaction details")

	var details details.Response

	err := a.timelineDetailsClient.Details(transaction.Action.Payload, &details)
	if err != nil {
		return fmt.Errorf("could not fetch transaction details: %w", err)
	}

	eventType, err := a.eventTypeResolver.Resolve(transaction)
	if err != nil {
		return fmt.Errorf("could not resolve transaction even type: %w", err)
	}

	a.logger.WithFields(infoFields).Info("Processing transaction details")

	if err := a.transactionProcessor.Process(eventType, details); err != nil {
		return fmt.Errorf("could not process transaction: %w", err)
	}

	a.logger.WithFields(infoFields).Info("Transaction processed")

	return nil
}

func (a App) GetActivityLog() ([]activitylog.ResponseItem, error) {
	a.logger.Info("Downloading activity entries")

	var entries []activitylog.ResponseItem

	err := a.activityLogClient.List(&entries)
	if err != nil {
		return entries, fmt.Errorf("could not fetch activity entries: %w", err)
	}

	slices.Reverse(entries)

	a.logger.Infof("%d activity entries downloaded", len(entries))

	return entries, nil
}
