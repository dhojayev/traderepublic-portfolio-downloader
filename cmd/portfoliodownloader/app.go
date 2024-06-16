package portfoliodownloader

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/activity"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

type App struct {
	activityHandler    activity.HandlerInterface
	transactionHandler transaction.HandlerInterface
	logger             *log.Logger
}

func NewApp(
	activityHandler activity.HandlerInterface,
	transactionHandler transaction.HandlerInterface,
	logger *log.Logger,
) App {
	return App{
		activityHandler:    activityHandler,
		transactionHandler: transactionHandler,
		logger:             logger,
	}
}

func (a App) Run() error {
	defer auth.SessionRefreshTicker.Stop()

	if err := a.transactionHandler.Handle(); err != nil {
		return fmt.Errorf("transaction handler errors: %w", err)
	}

	if err := a.activityHandler.Handle(); err != nil {
		return fmt.Errorf("activity handler error: %w", err)
	}

	return nil
}
