package timelinedetails

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	gocache "github.com/patrickmn/go-cache"
)

type Normalizer struct {
	builder *transaction.ModelBuilder
	cache   *gocache.Cache
}

func NewNormalizer(builder *transaction.ModelBuilder, cache *gocache.Cache) *Normalizer {
	return &Normalizer{
		builder: builder,
		cache:   cache,
	}
}

func (n *Normalizer) Normalize(data traderepublic.TimelineDetailsJson) error {
	parent, found := n.cache.Get(string(data.Id))
	if !found {
		return fmt.Errorf("timeline transaction %s not found in cache", data.Id)
	}

	item, ok := parent.(traderepublic.TimelineTransaction)
	if !ok {
		return fmt.Errorf("invalid timeline transaction in cache: %#v", parent)
	}

	if item.EventType != traderepublic.TimelineTransactionEventTypeINCOMINGTRANSFERDELEGATION {
		return nil
	}

	n.builder.
		WithID(string(data.Id)).
		WithType(transaction.TypeDeposit)

	var table traderepublic.TableSection

	err := data.Section(&table)
	if err != nil {
		return errors.New("failed to get table section")
	}

	for _, row := range table.Data {
		payment, ok := row.(traderepublic.PaymentRow)
		if !ok {
			continue
		}

		if payment.Title != "Gesamt" {
			continue
		}

		total, err := ParseFloatFromResponse(payment.Detail.Text)
		if err != nil {
			return fmt.Errorf("failed to convert string payment total amount to float: %w", err)
		}

		n.builder.WithCredit(total)
	}

	model := n.builder.Build()

	slog.Info("model built", "model", model)

	return nil
}
