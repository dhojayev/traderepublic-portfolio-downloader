package transaction

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

var ErrIgnoredTransactionReceived = errors.New("ignore transaction received")

type DataMapper struct {
}

func NewDataMapper() *DataMapper {
	return &DataMapper{}
}

func (m *DataMapper) Map(trnType TransactionType, details traderepublic.TimelineDetailsJson) (Model, error) {
	if trnType == TypeIgnored {
		return Model{}, fmt.Errorf("%w: %s", ErrIgnoredTransactionReceived, details.Id)
	}

	header, err := details.SectionHeader()
	if err != nil {
		return Model{}, fmt.Errorf("failed to find header section: %w", err)
	}

	_, err = details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return Model{}, fmt.Errorf("failed to find overview section: %w", err)
	}

	model := Model{
		ID:     string(details.Id),
		Status: string(header.Data.Status),
		Type:   trnType,
	}

	timestamp, err := ParseTimestamp(header.Data.Timestamp)
	if err != nil {
		return model, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	model.Timestamp = CSVDateTime{Time: timestamp}

	switch trnType {
	case TypeSavingsplan, TypeBuyOrder, TypeSellOrder, TypeRoundUp, TypeSaveback:
		model.ISIN = header.Action.Payload

		trn, err := details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return model, fmt.Errorf("failed to find transaction section: %w", err)
		}

		shares, err := trn.FindData(traderepublic.DataShares)
		if err != nil {
			return model, fmt.Errorf("failed to find shares data: %w", err)
		}

		model.Shares, err = strconv.ParseFloat(shares.Detail.Text, 64)
		if err != nil {
			return model, fmt.Errorf("failed to parse float from shares: %w", err)
		}

	case TypeDividendsIncome:
		model.ISIN, err = ExtractInstrumentISINFromIcon(header.Data.Icon)
		if err != nil {
			return model, fmt.Errorf("failed to extract ISIN from icon: %w", err)
		}
	}

	return model, nil
}
