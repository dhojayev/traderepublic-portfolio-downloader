package timelinedetails

import (
	"fmt"

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

func (n *Normalizer) Normalize(details traderepublic.TimelineDetailsJson) (transaction.Model, error) {
	// parent, found := n.cache.Get(string(details.Id))
	// if !found {
	// 	return fmt.Errorf("timeline transaction %s not found in cache", details.Id)
	// }

	// item, ok := parent.(traderepublic.TimelineTransaction)
	// if !ok {
	// 	return fmt.Errorf("invalid timeline transaction in cache: %#v", parent)
	// }

	var model transaction.Model

	n.builder.
		SetID(string(details.Id))

	header, err := details.SectionHeader()
	if err != nil {
		return model, fmt.Errorf("failed to find header section: %w", err)
	}

	timestamp, err := transaction.ParseTimestamp(header.Data.Timestamp)
	if err != nil {
		return model, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	n.builder.
		SetStatus(string(header.Data.Status)).
		SetTimestamp(transaction.CSVDateTime{Time: timestamp})

	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return model, fmt.Errorf("failed to find overview section: %w", err)
	}

	asset, err := overview.FindData(traderepublic.DataAsset)
	if err == nil {
		n.builder.SetAssetName(asset.Detail.Text)
	}

	return n.builder.Build(), nil
}
