package portfoliodownloader

import (
	"errors"
	"fmt"
	"slices"

	log "github.com/sirupsen/logrus"

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
	logger                *log.Logger
}

func NewApp(
	transactionsClient transactions.ClientInterface,
	eventTypeResolver transactions.EventTypeResolverInterface,
	timelineDetailsClient details.ClientInterface,
	transactionProcessor transaction.ProcessorInterface,
	logger *log.Logger,
) App {
	return App{
		transactionsClient:    transactionsClient,
		eventTypeResolver:     eventTypeResolver,
		timelineDetailsClient: timelineDetailsClient,
		transactionProcessor:  transactionProcessor,
		logger:                logger,
	}
}

func (a App) Run() error {
	counter := 0

	responses, err := a.GetTimelineTransactions()
	if err != nil {
		return err
	}

	for _, response := range responses {
		if response.Action.Payload == "" {
			continue
		}

		infoFields := log.Fields{"id": response.ID}

		err := a.ProcessTransaction(response)

		switch {
		case err == nil:
			counter++
		case errors.Is(err, websocket.ErrErrorStateReceived):
			a.logger.WithFields(infoFields).Error(err)

			continue
		case errors.Is(err, transactions.ErrUnsupportedEventType):
			a.logger.WithFields(infoFields).Warn("Unsupported transaction skipped")

			continue
		case errors.Is(err, transaction.ErrUnsupportedType):
			a.logger.WithFields(infoFields).Warn("Unsupported transaction skipped")

			continue
		}
	}

	skippedCount := len(responses) - counter

	a.logger.Infof("Completed: %d, skipped: %d", counter, skippedCount)

	return nil
}

func (a App) GetTimelineTransactions() ([]transactions.ResponseItem, error) {
	a.logger.Info("Downloading transactions")

	transactionResponses, err := a.transactionsClient.Get()
	if err != nil {
		return transactionResponses, fmt.Errorf("could not fetch transactions: %w", err)
	}

	slices.Reverse(transactionResponses)

	a.logger.Infof("%d transactions downloaded", len(transactionResponses))

	return transactionResponses, nil
}

func (a App) ProcessTransaction(response transactions.ResponseItem) error {
	infoFields := log.Fields{"id": response.ID}

	a.logger.WithFields(infoFields).Info("Fetching transaction details")

	transactionDetails, err := a.timelineDetailsClient.Get(response.Action.Payload)
	if err != nil {
		return fmt.Errorf("could not fetch transaction details: %w", err)
	}

	eventType, err := a.eventTypeResolver.Resolve(response)
	if err != nil {
		return fmt.Errorf("could not resolve transaction even type: %w", err)
	}

	a.logger.WithFields(infoFields).Info("Processing transaction details")

	if err := a.transactionProcessor.Process(eventType, transactionDetails); err != nil {
		return fmt.Errorf("could process transaction: %w", err)
	}

	a.logger.WithFields(infoFields).Info("Transaction processed")

	return nil
}
