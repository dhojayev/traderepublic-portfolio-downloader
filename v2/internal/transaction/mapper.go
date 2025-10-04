package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	gocache "github.com/patrickmn/go-cache"
)

var ErrTransactionWithoutTypeReceived = errors.New("transaction without type received")

type DataMapper struct {
	cache *gocache.Cache
}

func NewDataMapper(cache *gocache.Cache) *DataMapper {
	return &DataMapper{
		cache: cache,
	}
}

func (m *DataMapper) Map(details traderepublic.TimelineDetailsJson, model *Model) error {
	if model.Type == nil {
		return fmt.Errorf("%w: %s", ErrTransactionWithoutTypeReceived, details.Id)
	}

	model.ID = model.Type.FindID(details)
	status, err := model.Type.FindStatus(details)
	if err != nil {
		return fmt.Errorf("failed to find status in details: %w", err)
	}

	model.Status = status

	timestampStr, err := model.Type.FindTimestamp(details)
	if err != nil {
		return fmt.Errorf("failed to find timestamp in details: %w", err)
	}

	timestamp, err := ParseTimestamp(timestampStr)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	model.Timestamp = CSVDateTime{Time: timestamp}

	isin, err := model.Type.FindISIN(details)
	if err != nil {
		return fmt.Errorf("failed to find ISIN in details: %w", err)
	}

	model.ISIN = isin

	sharesStr, err := model.Type.FindShares(details)
	if err != nil {
		return fmt.Errorf("failed to find shares in details: %w", err)
	}

	shares, err := ParseFloatFromResponse(sharesStr)
	if err != nil {
		return fmt.Errorf("failed to parse float from shares: %w", err)
	}

	model.Shares = shares

	shaePriceStr, err := model.Type.FindSharePrice(details)
	if err != nil {
		return fmt.Errorf("failed to find share price in details: %w", err)
	}

	sharePrice, err := ParseFloatFromResponse(shaePriceStr)
	if err != nil {
		return fmt.Errorf("failed to parse float from share price: %w", err)
	}

	model.SharePrice = sharePrice

	feeStr, err := model.Type.FindFee(details)
	if err != nil {
		return fmt.Errorf("failed to find fee data: %w", err)
	}

	zeroFee := float64(0)
	model.Fee = &zeroFee

	if feeStr != "Free" {
		fee, err := ParseFloatFromResponse(feeStr)
		if err != nil {
			return fmt.Errorf("failed to parse float from fee: %w", err)
		}

		model.Fee = &fee
	}

	totalStr, err := model.Type.FindTotal(details)
	if err != nil {
		return fmt.Errorf("failed to find total data: %w", err)
	}

	total, err := ParseFloatFromResponse(totalStr)
	if err != nil {
		return fmt.Errorf("failed to parse float from total: %w", err)
	}

	model.Debit = total

	if isin != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		instr, err := m.getInstrument(ctx, isin)
		if err != nil {
			return fmt.Errorf("failed to get instrument: %w", err)
		}

		model.AssetName = *instr.ShortName
		model.AssetType = string(instr.TypeId)
	}

	return nil

	// header, err := details.SectionHeader()
	// if err != nil {
	// 	return Model{}, fmt.Errorf("failed to find header section: %w", err)
	// }

	// overview, err := details.FindSection(traderepublic.SectionOverview)
	// if err != nil {
	// 	return Model{}, fmt.Errorf("failed to find overview section: %w", err)
	// }

	// model := Model{
	// 	ID:     string(details.Id),
	// 	Status: string(header.Data.Status),
	// }

	// timestamp, err := ParseTimestamp(header.Data.Timestamp)
	// if err != nil {
	// 	return model, fmt.Errorf("failed to parse timestamp: %w", err)
	// }

	// model.Timestamp = CSVDateTime{Time: timestamp}

	// switch trnType {
	// case TypeBuyOrder, TypeSellOrder, TypeRoundUp, TypeSaveback:
	// 	model.ISIN = header.Action.Payload

	// 	trn, err := details.FindSection(traderepublic.SectionTransaction)
	// 	if err != nil {
	// 		return model, fmt.Errorf("failed to find transaction section: %w", err)
	// 	}

	// 	shares, err := trn.FindData(traderepublic.DataShares)
	// 	if err != nil {
	// 		return model, fmt.Errorf("failed to find shares data: %w", err)
	// 	}

	// 	model.Shares, err = strconv.ParseFloat(shares.Detail.Text, 64)
	// 	if err != nil {
	// 		return model, fmt.Errorf("failed to parse float from shares: %w", err)
	// 	}

	// case TypeSavingsplan:
	// 	model.ISIN = header.Action.Payload

	// 	var sharesVal string

	// 	trn, err := overview.FindData(traderepublic.DataTransaction)
	// 	if err != nil {
	// 		trnSection, err := details.FindSection(traderepublic.SectionTransaction)
	// 		if err != nil {

	// 		}
	// 		if err == nil {
	// 			shares, err := trnSection.FindData(traderepublic.DataShares)
	// 			if err == nil {
	// 				sharesVal = shares.Detail.Text

	// 				break
	// 			}
	// 		}

	// 		return model, fmt.Errorf("failed to find transaction data: %w", err)
	// 	}

	// 	if trn.Detail.DisplayValue != nil {
	// 		sharesVal = *trn.Detail.DisplayValue.Prefix
	// 	}

	// 	model.Shares, err = ParseFloatFromResponse(sharesVal)
	// 	if err != nil {
	// 		return model, fmt.Errorf("failed to parse shares data: %w", err)
	// 	}
	// case TypeDividendsIncome:
	// 	model.ISIN, err = ExtractInstrumentISINFromIcon(header.Data.Icon)
	// 	if err != nil {
	// 		return model, fmt.Errorf("failed to extract ISIN from icon: %w", err)
	// 	}
	// }

	// return model, nil
}

func (m *DataMapper) getInstrument(ctx context.Context, isin string) (traderepublic.InstrumentJson, error) {

	for {
		select {
		case <-ctx.Done():
			return traderepublic.InstrumentJson{}, fmt.Errorf("context timeout: %w", ctx.Err())

		default:
			entry, found := m.cache.Get(isin)
			if found {
				instr, ok := entry.(traderepublic.InstrumentJson)
				if !ok {
					continue
				}

				return instr, nil
			}
		}
	}
}
