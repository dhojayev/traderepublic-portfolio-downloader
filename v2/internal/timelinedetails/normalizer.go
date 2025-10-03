package timelinedetails

import (
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type Normalizer struct {
	builder *transaction.ModelBuilder
}

func NewNormalizer(builder *transaction.ModelBuilder) *Normalizer {
	return &Normalizer{
		builder: builder,
	}
}

func (n *Normalizer) Normalize(details traderepublic.TimelineDetailsJson) (transaction.Model, error) {
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
