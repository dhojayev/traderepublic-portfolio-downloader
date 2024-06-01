package portfoliodownloader

import (
	"errors"
	"fmt"
	"slices"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

type App struct {
	transactionsClient    transactions.Client
	eventTypeResolver     transactions.EventTypeResolverInterface
	timelineDetailsClient details.Client
	transactionProcessor  transaction.Processor
	logger                *log.Logger
}

func NewApp(
	transactionsClient transactions.Client,
	eventTypeResolver transactions.EventTypeResolverInterface,
	timelineDetailsClient details.Client,
	transactionProcessor transaction.Processor,
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
	a.logger.Info("Downloading transactions")

	transactionResponses, err := a.transactionsClient.Get()
	if err != nil {
		return fmt.Errorf("could not fetch transactions: %w", err)
	}

	a.logger.Infof("%d transactions downloaded", len(transactionResponses))

	slices.Reverse(transactionResponses)

	for _, transactionResponse := range transactionResponses {
		if transactionResponse.Action.Payload == "" {
			continue
		}

		id := transactionResponse.Action.Payload
		infoFields := log.Fields{"id": id}

		a.logger.WithFields(infoFields).Info("Fetching transaction details")

		transactionDetails, err := a.timelineDetailsClient.Get(id)
		if err != nil {
			return fmt.Errorf("could not fetch transaction details: %w", err)
		}

		eventType, err := a.eventTypeResolver.Resolve(transactionResponse)
		if err != nil {
			if errors.Is(err, transactions.ErrUnsupportedEventType) {
				a.logger.WithFields(infoFields).Info("Unsupported transaction skipped")

				continue
			}

			return fmt.Errorf("could not resolve transaction even type: %w", err)
		}

		a.logger.WithFields(infoFields).Info("Processing transaction details")

		if err := a.transactionProcessor.Process(eventType, transactionDetails); err != nil {
			if errors.Is(err, transaction.ErrUnsupportedType) {
				a.logger.WithFields(infoFields).Info("Unsupported transaction skipped")

				continue
			}

			return fmt.Errorf("could process transaction: %w", err)
		}

		a.logger.WithFields(infoFields).Info("Transaction processed")
	}

	a.logger.Info("All data processed")

	return nil
}
