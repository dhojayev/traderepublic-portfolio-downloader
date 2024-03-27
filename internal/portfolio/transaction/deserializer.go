package transaction

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
)

var ErrUnsupportedResponse = errors.New("unsupported response")

type DetailsDeserializerInterface interface {
	Deserialize(response details.Response) (any, error)
}

type DetailsDeserializer struct {
	resolver TypeResolver
	logger   *log.Logger
}

func NewDetailsDeserializer(resolver TypeResolver, logger *log.Logger) DetailsDeserializer {
	return DetailsDeserializer{
		resolver: resolver,
		logger:   logger,
	}
}

func (d DetailsDeserializer) Deserialize(response details.Response) (any, error) {
	resolvedType, err := d.resolver.Resolve(response)
	if err != nil {
		return Purchase{}, fmt.Errorf("resolver error: %w", err)
	}

	switch resolvedType {
	case TypeSaleTransaction:
		return d.DeserializeSaleTransaction(response)
	case TypePurchaseTransaction:
		return d.DeserializePurchase(response)
	case TypeDividendPayoutTransaction:
		return d.DeserializeDividendPayout(response)
	case TypeRoundUpTransaction:
		return d.DeserializeBenefit(TypeRoundUp, response)
	case TypeSavebackTransaction:
		return d.DeserializeBenefit(TypeSaveback, response)
	case TypeUnsupported, TypeCardPaymentTransaction:
		return Purchase{}, ErrUnsupportedResponse
	default:
		return Purchase{}, ErrUnsupportedResponse
	}
}

func (d DetailsDeserializer) DeserializePurchase(response details.Response) (Purchase, error) {
	transaction, err := d.DeserializeBaseTransaction(response)
	if err != nil {
		return Purchase{}, err
	}

	asset, err := d.DeserializeAsset(response)
	if err != nil {
		return Purchase{}, err
	}

	monetaryValues, err := d.DeserializeMonetaryValues(response)
	if err != nil {
		return Purchase{}, err
	}

	documents, err := d.DeserializeDocuments(response)
	if err != nil {
		return Purchase{}, err
	}

	return NewPurchase(transaction, asset, monetaryValues, documents), nil
}

func (d DetailsDeserializer) DeserializeSaleTransaction(response details.Response) (Sale, error) {
	purchase, err := d.DeserializePurchase(response)
	if err != nil {
		return Sale{}, err
	}

	sale := NewSale(0, 0, purchase)

	performance, err := response.PerformanceSection()
	if err != nil {
		if errors.Is(err, details.ErrSectionNotFound) {
			return sale, nil
		}

		return sale, fmt.Errorf("could get performance section: %w", err)
	}

	yield, err := performance.Yield()
	if err != nil {
		return sale, fmt.Errorf("could get yield: %w", err)
	}

	yieldParsed, err := ParseFloatWithComma(yield.Detail.Text, yield.Detail.IsTrendNegative())
	if err != nil {
		return sale, fmt.Errorf("could not parse float value from yield: %w", err)
	}

	sale.Yield = yieldParsed

	profit, err := performance.Profit()
	if err != nil {
		return sale, fmt.Errorf("could get profit: %w", err)
	}

	profitParsed, err := ParseFloatWithComma(profit.Detail.Text, profit.Detail.IsTrendNegative())
	if err != nil {
		return sale, fmt.Errorf("could not parse float value from profit: %w", err)
	}

	sale.Profit = profitParsed

	return sale, nil
}

func (d DetailsDeserializer) DeserializeDividendPayout(response details.Response) (DividendPayout, error) {
	purchase, err := d.DeserializePurchase(response)
	if err != nil {
		return DividendPayout{}, err
	}

	return NewDividendPayout(NewSale(0, 0, purchase)), nil
}

func (d DetailsDeserializer) DeserializeBenefit(benefitType string, response details.Response) (Benefit, error) {
	purchase, err := d.DeserializePurchase(response)
	if err != nil {
		return Benefit{}, err
	}

	return NewBenefit(benefitType, purchase), nil
}

func (d DetailsDeserializer) DeserializeBaseTransaction(response details.Response) (BaseTransaction, error) {
	header, err := response.HeaderSection()
	if err != nil {
		return BaseTransaction{}, fmt.Errorf("could not get transaction header %w", err)
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05-0700", header.Data.Timestamp)
	if err != nil {
		return BaseTransaction{}, fmt.Errorf("could not parse timestamp: %w", err)
	}

	return NewBaseTransaction(response.ID, header.Data.Status, timestamp), nil
}

func (d DetailsDeserializer) DeserializeAsset(response details.Response) (Asset, error) {
	overview, err := response.OverviewSection()
	if err != nil {
		return Asset{}, fmt.Errorf("error getting overview: %w", err)
	}

	asset, err := overview.Asset()
	if err != nil {
		return Asset{}, fmt.Errorf("error getting asset: %w", err)
	}

	header, err := response.HeaderSection()
	if err != nil {
		return Asset{}, fmt.Errorf("error getting transaction header: %w", err)
	}

	instrument, _ := header.Action.Payload.(string)
	if instrument == "" {
		instrument, _ = ExtractInstrumentNameFromIcon(header.Data.Icon)
	}

	transactionSection, err := response.TransactionSection()
	if err != nil {
		return Asset{}, fmt.Errorf("could not get transaction section: %w", err)
	}

	shares, err := transactionSection.Shares()
	if err != nil {
		return Asset{}, fmt.Errorf("could not get shares: %w", err)
	}

	sharesParsed, err := ParseFloatWithComma(shares.Detail.Text, shares.Detail.IsTrendNegative())
	if err != nil {
		return Asset{}, fmt.Errorf("could not parse float value from shares: %w", err)
	}

	if shares.HasSharesWithPeriod() {
		sharesParsed, err = ParseFloatWithPeriod(shares.Detail.Text)
		if err != nil {
			return Asset{}, fmt.Errorf("could not parse float value from shares: %w", err)
		}
	}

	return NewAsset(instrument, asset.Detail.Text, sharesParsed), nil
}

func (d DetailsDeserializer) DeserializeMonetaryValues(response details.Response) (MonetaryValues, error) {
	transactionSection, err := response.TransactionSection()
	if err != nil {
		return MonetaryValues{}, fmt.Errorf("could not get transaction section: %w", err)
	}

	rate, err := transactionSection.Rate()
	if err != nil {
		return MonetaryValues{}, fmt.Errorf("could not get rate: %w", err)
	}

	rateParsed, err := ParseFloatWithComma(rate.Detail.Text, rate.Detail.IsTrendNegative())
	if err != nil {
		return MonetaryValues{}, fmt.Errorf("could not parse float value from rate: %w", err)
	}

	commission, err := transactionSection.Commission()
	if err != nil {
		if !errors.Is(err, details.ErrSectionDataEntryNotFound) {
			return MonetaryValues{}, fmt.Errorf("could not get commission: %w", err)
		}
	}

	commissionParsed, err := ParseFloatWithComma(commission.Detail.Text, commission.Detail.IsTrendNegative())
	if err != nil {
		if !errors.Is(err, ErrNoMatch) {
			return MonetaryValues{}, fmt.Errorf("could not parse float value from commission: %w", err)
		}

		commissionParsed = 0
	}

	total, err := transactionSection.Total()
	if err != nil {
		return MonetaryValues{}, fmt.Errorf("could not get total: %w", err)
	}

	totalParsed, err := ParseFloatWithComma(total.Detail.Text, total.Detail.IsTrendNegative())
	if err != nil {
		return MonetaryValues{}, fmt.Errorf("could not parse float value from total: %w", err)
	}

	return NewMonetaryValues(rateParsed, commissionParsed, totalParsed), nil
}

func (d DetailsDeserializer) DeserializeDocuments(response details.Response) ([]Document, error) {
	documents := make([]Document, 0)

	documentsSection, err := response.DocumentsSection()
	if err != nil {
		return documents, fmt.Errorf("could not get documents section: %w", err)
	}

	for _, document := range documentsSection.Data {
		url, ok := document.Action.Payload.(string)
		if !ok {
			continue
		}

		documents = append(documents, Document{
			ID:    document.ID,
			URL:   url,
			Date:  document.Detail,
			Title: document.Title,
		})
	}

	return documents, nil
}
