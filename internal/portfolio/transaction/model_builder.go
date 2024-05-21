package transaction

import (
	"errors"
	"fmt"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	log "github.com/sirupsen/logrus"
)

type ModelBuilderFactoryInterface interface {
	Create(response details.Response) (ModelBuilderInterface, error)
}

type ModelBuilderFactory struct {
	resolver details.TypeResolverInterface
	logger   *log.Logger
}

func NewModelBuilderFactory(resolver details.TypeResolverInterface, logger *log.Logger) ModelBuilderFactory {
	return ModelBuilderFactory{
		resolver: resolver,
		logger:   logger,
	}
}

//nolint:ireturn
func (f ModelBuilderFactory) Create(response details.Response) (ModelBuilderInterface, error) {
	responseType, err := f.resolver.Resolve(response)
	if err != nil {
		return nil, fmt.Errorf("resolver errors: %w", err)
	}

	baseBuilder := BaseModelBuilder{response: response, logger: f.logger}

	switch responseType {
	case details.TypePurchaseTransaction:
		return PurchaseBuilder{baseBuilder}, nil
	case details.TypeSaleTransaction:
		return SaleBuilder{PurchaseBuilder{baseBuilder}}, nil
	case details.TypeDividendPayoutTransaction:
		return DividendPayoutBuilder{SaleBuilder{PurchaseBuilder{baseBuilder}}}, nil
	case details.TypeRoundUpTransaction:
		return RoundUpBuilder{PurchaseBuilder{baseBuilder}}, nil
	case details.TypeSavebackTransaction:
		return SavebackBuilder{PurchaseBuilder{baseBuilder}}, nil
	case
		details.TypeUnsupported,
		details.TypeCardPaymentTransaction,
		details.TypeDepositTransaction,
		details.TypeDepositInterestReceivedTransaction:
		return nil, ErrUnsupportedType
	}

	return nil, ErrUnsupportedType
}

type ModelBuilderInterface interface {
	Build() (Model, error)
}

type BaseModelBuilder struct {
	response details.Response
	logger   *log.Logger
}

func (b BaseModelBuilder) ExtractStatus() (string, error) {
	header, err := b.response.HeaderSection()
	if err != nil {
		return "", fmt.Errorf("could not get header section: %w", err)
	}

	return header.Data.Status, nil
}

func (b BaseModelBuilder) ExtractInstrumentIcon() (string, error) {
	header, err := b.response.HeaderSection()
	if err != nil {
		return "", fmt.Errorf("could not get header section: %w", err)
	}

	return header.Data.Icon, nil
}

func (b BaseModelBuilder) ExtractTimestamp() (time.Time, error) {
	header, err := b.response.HeaderSection()
	if err != nil {
		return time.Time{}, fmt.Errorf("could not get header section: %w", err)
	}

	timestamp, err := time.Parse(internal.DefaultTimeFormat, header.Data.Timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse header section timestamp: %w", err)
	}

	return timestamp, nil
}

func (b BaseModelBuilder) ExtractInstrumentISIN() (string, error) {
	header, err := b.response.HeaderSection()
	if err != nil {
		return "", fmt.Errorf("could not get header section: %w", err)
	}

	isin, _ := header.Action.Payload.(string)
	if isin == "" {
		isin, _ = ExtractInstrumentNameFromIcon(header.Data.Icon)
	}

	return isin, nil
}

func (b BaseModelBuilder) ExtractInstrumentName() (string, error) {
	overview, err := b.response.OverviewSection()
	if err != nil {
		return "", fmt.Errorf("could not get overview section: %w", err)
	}

	asset, err := overview.Asset()
	if err == nil {
		return asset.Detail.Text, nil
	}

	if !errors.Is(err, details.ErrSectionDataEntryNotFound) {
		return "", fmt.Errorf("could not get overview section asset: %w", err)
	}

	underlyingAsset, err := overview.UnderlyingAsset()
	if err != nil {
		return "", fmt.Errorf("could not get overview section underlying asset: %w", err)
	}

	return underlyingAsset.Detail.Text, nil
}

func (b BaseModelBuilder) ExtractYield() (float64, error) {
	performance, err := b.response.PerformanceSection()
	if err != nil {
		return 0, fmt.Errorf("could not get performance section: %w", err)
	}

	yieldData, err := performance.Yield()
	if err != nil {
		b.logger.Debugf("could not get performance section yield: %s", err)
	}

	yield, err := ParseFloatWithComma(yieldData.Detail.Text, yieldData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse performance section yield to float: %s", err)
	}

	return yield, nil
}

func (b BaseModelBuilder) ExtractProfit() (float64, error) {
	performance, err := b.response.PerformanceSection()
	if err != nil {
		return 0, fmt.Errorf("could not get performance section: %w", err)
	}

	profitData, err := performance.Profit()
	if err != nil {
		b.logger.Debugf("could not get performance section profit: %s", err)
	}

	profit, err := ParseFloatWithComma(profitData.Detail.Text, profitData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse performance section profit to float: %s", err)
	}

	return profit, nil
}

func (b BaseModelBuilder) ExtractSharesAmount() (float64, error) {
	transactionSection, err := b.response.TransactionSection()
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	sharesData, err := transactionSection.Shares()
	if err != nil {
		b.logger.Debugf("could not get transaction section shares: %s", err)
	}

	shares, err := ParseFloatWithComma(sharesData.Detail.Text, sharesData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse transaction section shares to float: %s", err)
	}

	if sharesData.HasSharesWithPeriod() {
		shares, err = ParseFloatWithPeriod(sharesData.Detail.Text)
		if err != nil {
			b.logger.Debugf("could not parse transaction section shares to float: %s", err)
		}
	}

	return shares, nil
}

func (b BaseModelBuilder) ExtractRateValue() (float64, error) {
	transactionSection, err := b.response.TransactionSection()
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	rateData, err := transactionSection.Rate()
	if err != nil {
		b.logger.Debugf("could not get transaction section: %s", err)
	}

	rate, err := ParseFloatWithComma(rateData.Detail.Text, rateData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse transaction section rate to float: %s", err)
	}

	return rate, nil
}

func (b BaseModelBuilder) ExtractCommissionAmount() (float64, error) {
	transactionSection, err := b.response.TransactionSection()
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	commissionData, err := transactionSection.Commission()
	if err != nil {
		if !errors.Is(err, details.ErrSectionDataEntryNotFound) {
			b.logger.Debugf("could not get transaction section commission: %s", err)
		}
	}

	commission, err := ParseFloatWithComma(commissionData.Detail.Text, commissionData.Detail.IsTrendNegative())
	if err != nil {
		if !errors.Is(err, ErrNoMatch) {
			b.logger.Debugf("could not parse transaction section commission to float: %s", err)
		}

		commission = 0
	}

	return commission, nil
}

func (b BaseModelBuilder) ExtractTotalAmount() (float64, error) {
	transactionSection, err := b.response.TransactionSection()
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	totalData, err := transactionSection.Total()
	if err != nil {
		b.logger.Debugf("could not get transaction section total: %s", err)
	}

	total, err := ParseFloatWithComma(totalData.Detail.Text, totalData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse transaction section total to float: %s", err)
	}

	return total, nil
}

func (b BaseModelBuilder) BuildDocuments() ([]document.Model, error) {
	documents := make([]document.Model, 0)

	documentsSection, err := b.response.DocumentsSection()
	if err != nil {
		return documents, fmt.Errorf("could not get documents section: %w", err)
	}

	for _, doc := range documentsSection.Data {
		url, ok := doc.Action.Payload.(string)
		if !ok {
			continue
		}

		documents = append(documents, document.NewModel(doc.ID, url, doc.Detail, doc.Title))
	}

	return documents, nil
}

type PurchaseBuilder struct {
	BaseModelBuilder
}

func (b PurchaseBuilder) Build() (Model, error) {
	var err error

	model := Model{
		UUID: b.response.ID,
		Type: TypePurchase,
	}

	model.Status, err = b.ExtractStatus()
	if err != nil {
		return model, err
	}

	model.Timestamp, err = b.ExtractTimestamp()
	if err != nil {
		return model, err
	}

	model.Instrument.ISIN, err = b.ExtractInstrumentISIN()
	if err != nil {
		return model, err
	}

	model.Instrument.Name, err = b.ExtractInstrumentName()
	if err != nil {
		return model, err
	}

	model.Shares, err = b.ExtractSharesAmount()
	if err != nil {
		return model, err
	}

	model.Rate, err = b.ExtractRateValue()
	if err != nil {
		return model, err
	}

	model.Commission, err = b.ExtractCommissionAmount()
	if err != nil {
		return model, err
	}

	model.Total, err = b.ExtractTotalAmount()
	if err != nil {
		return model, err
	}

	model.Instrument.Icon, _ = b.ExtractInstrumentIcon()
	model.Documents, _ = b.BuildDocuments()

	return model, err
}

type SaleBuilder struct {
	PurchaseBuilder
}

func (b SaleBuilder) Build() (Model, error) {
	model, err := b.PurchaseBuilder.Build()
	if err != nil {
		return model, err
	}

	model.Type = TypeSale

	model.Yield, err = b.ExtractYield()
	if err != nil {
		return model, err
	}

	model.Profit, err = b.ExtractProfit()
	if err != nil {
		return model, err
	}

	return model, nil
}

type RoundUpBuilder struct {
	PurchaseBuilder
}

func (b RoundUpBuilder) Build() (Model, error) {
	model, err := b.PurchaseBuilder.Build()
	if err != nil {
		return model, err
	}

	model.Type = TypeRoundUp

	return model, nil
}

type SavebackBuilder struct {
	PurchaseBuilder
}

func (b SavebackBuilder) Build() (Model, error) {
	model, err := b.PurchaseBuilder.Build()
	if err != nil {
		return model, err
	}

	model.Type = TypeSaveback

	return model, nil
}

type DividendPayoutBuilder struct {
	SaleBuilder
}

func (b DividendPayoutBuilder) Build() (Model, error) {
	model, err := b.SaleBuilder.Build()
	if err != nil {
		return model, err
	}

	model.Type = TypeDividendPayout

	return model, nil
}
