package activity

import (
	"fmt"
	"slices"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
)

type HandlerInterface interface {
	Handle() error
}

type Handler struct {
	listClient    activitylog.ClientInterface
	detailsClient details.ClientInterface
	normalizer    details.ResponseNormalizerInterface
	processor     ProcessorInterface
	logger        *log.Logger
}

func NewHandler(
	listClient activitylog.ClientInterface,
	detailsClient details.ClientInterface,
	normalizer details.ResponseNormalizerInterface,
	processor ProcessorInterface,
	logger *log.Logger,
) Handler {
	return Handler{
		listClient:    listClient,
		detailsClient: detailsClient,
		normalizer:    normalizer,
		processor:     processor,
		logger:        logger,
	}
}

func (h Handler) Handle() error {
	counter := internal.NewOperationCounter()

	entries, err := h.GetActivityLog()
	if err != nil {
		return fmt.Errorf("could not fetch activity log entries: %w", err)
	}

	for _, entry := range entries {
		if !entry.Action.HasDetails() {
			counter.Skipped().Add(1)

			continue
		}

		infoFields := log.Fields{"id": entry.ID}

		h.logger.WithFields(infoFields).Info("Fetching activity log entry details")

		var detailsEntry details.Response

		err := h.detailsClient.Details(entry.Action.PayloadStr(), &detailsEntry)
		if err != nil {
			return fmt.Errorf("could not fetch activity log entry details: %w", err)
		}

		normalizedResponse, _ := h.normalizer.Normalize(detailsEntry)
		if normalizedResponse.Documents == nil {
			counter.Skipped().Add(1)

			continue
		}

		if err := h.processor.Process(normalizedResponse); err != nil {
			return fmt.Errorf("activity processor error: %w", err)
		}

		counter.Processed().Add(1)
	}

	h.logger.Infof(
		"Activity log entries completed: %d, skipped: %d",
		counter.Processed().Load(),
		counter.Skipped().Load(),
	)

	return nil
}

func (h Handler) GetActivityLog() ([]activitylog.ResponseItem, error) {
	h.logger.Info("Downloading activity log entries")

	var entries []activitylog.ResponseItem

	err := h.listClient.List(&entries)
	if err != nil {
		return entries, fmt.Errorf("could not fetch activity log entries: %w", err)
	}

	slices.Reverse(entries)

	h.logger.Infof("%d activity log entries downloaded", len(entries))

	return entries, nil
}
