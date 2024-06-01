package transaction

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
)

var ErrUnsupportedType = errors.New("unsupported response")

type ModelBuilderFactoryInterface interface {
	Create(eventType transactions.EventType, response details.Response) (ModelBuilderInterface, error)
}

type ModelBuilderFactory struct {
	resolver         details.TypeResolverInterface
	documentsBuilder document.ModelBuilderInterface
	logger           *log.Logger
}

func NewModelBuilderFactory(
	resolver details.TypeResolverInterface,
	documentsBuilder document.ModelBuilderInterface,
	logger *log.Logger,
) ModelBuilderFactory {
	return ModelBuilderFactory{
		resolver:         resolver,
		documentsBuilder: documentsBuilder,
		logger:           logger,
	}
}

//nolint:ireturn,cyclop
func (f ModelBuilderFactory) Create(
	eventType transactions.EventType,
	response details.Response,
) (ModelBuilderInterface, error) {
	responseType, err := f.resolver.Resolve(eventType, response)
	if err != nil {
		if errors.Is(err, details.ErrUnsupportedResponse) {
			return nil, ErrUnsupportedType
		}

		return nil, fmt.Errorf("resolver error: %w", err)
	}

	baseBuilder := NewBaseModelBuilder(response, f.documentsBuilder, f.logger)
	purchaseBuilder := NewPurchaseBuilder(baseBuilder)

	switch responseType {
	case details.TypePurchaseTransaction:
		return purchaseBuilder, nil
	case details.TypeSaleTransaction:
		return NewSaleBuilder(purchaseBuilder), nil
	case details.TypeDividendPayoutTransaction:
		return NewDividendPayoutBuilder(purchaseBuilder), nil
	case details.TypeRoundUpTransaction:
		return NewRoundUpBuilder(purchaseBuilder), nil
	case details.TypeSavebackTransaction:
		return NewSavebackBuilder(purchaseBuilder), nil
	case details.TypeDepositTransaction:
		return NewDepositBuilder(baseBuilder), nil
	case details.TypeWithdrawalTransaction:
		return NewWithdrawBuilder(NewDepositBuilder(baseBuilder)), nil
	case
		details.TypeUnsupported,
		details.TypeCardPaymentTransaction,
		details.TypeInterestPayoutTransaction:
		return nil, ErrUnsupportedType
	}

	return nil, ErrUnsupportedType
}

type ModelBuilderInterface interface {
	Build() (Model, error)
}

type BaseModelBuilder struct {
	response         details.Response
	documentsBuilder document.ModelBuilderInterface
	logger           *log.Logger
}

func NewBaseModelBuilder(
	response details.Response,
	documentsBuilder document.ModelBuilderInterface,
	logger *log.Logger,
) BaseModelBuilder {
	return BaseModelBuilder{
		response:         response,
		documentsBuilder: documentsBuilder,
		logger:           logger,
	}
}

func (b BaseModelBuilder) ExtractStatus() (string, error) {
	header, err := b.response.SectionTypeHeader()
	if err != nil {
		return "", fmt.Errorf("could not get header section: %w", err)
	}

	return header.Data.Status, nil
}

func (b BaseModelBuilder) ExtractInstrumentIcon() (string, error) {
	header, err := b.response.SectionTypeHeader()
	if err != nil {
		return "", fmt.Errorf("could not get header section: %w", err)
	}

	return header.Data.Icon, nil
}

func (b BaseModelBuilder) ExtractTimestamp() (time.Time, error) {
	header, err := b.response.SectionTypeHeader()
	if err != nil {
		return time.Time{}, fmt.Errorf("could not get header section: %w", err)
	}

	timestamp, err := time.Parse(details.ResponseTimeFormat, header.Data.Timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse header section timestamp: %w", err)
	}

	return timestamp, nil
}

func (b BaseModelBuilder) ExtractInstrumentISIN() (string, error) {
	header, err := b.response.SectionTypeHeader()
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
	tableSections, err := b.response.SectionsTypeTable()
	if err != nil {
		return "", fmt.Errorf("could not get table sections: %w", err)
	}

	overview, err := tableSections.FindByTitle(details.SectionTitleOverview)
	if err != nil {
		return "", fmt.Errorf("could not get overview section: %w", err)
	}

	asset, err := overview.GetDataByTitle(details.OverviewDataTitleAsset)
	if err == nil {
		return asset.Detail.Text, nil
	}

	if !errors.Is(err, details.ErrSectionDataTitleNotFound) {
		return "", fmt.Errorf("could not get overview section asset: %w", err)
	}

	underlyingAsset, err := overview.GetDataByTitle(details.OverviewDataTitleUnderlyingAsset)
	if err != nil {
		return "", fmt.Errorf("could not get overview section underlying asset: %w", err)
	}

	return underlyingAsset.Detail.Text, nil
}

func (b BaseModelBuilder) ExtractYield() (float64, error) {
	tableSections, err := b.response.SectionsTypeHorizontalTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	performance, err := tableSections.FindByTitle(details.SectionTitlePerformance)
	if err != nil {
		return 0, fmt.Errorf("could not get performance section: %w", err)
	}

	yieldData, err := performance.GetDataByTitle(details.PerformanceDataTitleYield)
	if err != nil {
		return 0, fmt.Errorf("could not get performance section yield: %w", err)
	}

	yield, err := ParseFloatWithComma(yieldData.Detail.Text, yieldData.Detail.Trend == details.TrendNegative)
	if err != nil {
		return 0, fmt.Errorf("could not parse performance section yield to float: %w", err)
	}

	return yield, nil
}

func (b BaseModelBuilder) ExtractProfitAndLoss() (float64, error) {
	tableSections, err := b.response.SectionsTypeHorizontalTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	performance, err := tableSections.FindByTitle(details.SectionTitlePerformance)
	if err != nil {
		return 0, fmt.Errorf("could not get performance section: %w", err)
	}

	var profitData details.ResponseSectionTypeTableData

	titles := []string{details.PerformanceDataTitleProfit, details.PerformanceDataTitleLoss}

	for _, title := range titles {
		profitData, err = performance.GetDataByTitle(title)
		if err == nil {
			break
		}
	}

	if err != nil {
		return 0, fmt.Errorf("could not get performance section profit (%s): %w", titles, err)
	}

	profit, err := ParseFloatWithComma(profitData.Detail.Text, profitData.Detail.Trend == details.TrendNegative)
	if err != nil {
		return 0, fmt.Errorf("could not parse performance section profit to float: %w", err)
	}

	return profit, nil
}

func (b BaseModelBuilder) ExtractSharesAmount() (float64, error) {
	tableSections, err := b.response.SectionsTypeTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	transactionSection, err := tableSections.FindByTitle(details.SectionTitleTransaction)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	var sharesData details.ResponseSectionTypeTableData

	titles := []string{
		details.TransactionDataTitleShares,
		details.TransactionDataTitleSharesAlt,
	}

	for _, title := range titles {
		sharesData, err = transactionSection.GetDataByTitle(title)
		if err == nil {
			break
		}
	}

	if err != nil {
		return 0, fmt.Errorf("could not get transaction section shares (%s): %w", titles, err)
	}

	shares, err := ParseFloatWithPeriod(sharesData.Detail.Text)
	if err != nil {
		shares, err = ParseFloatWithComma(sharesData.Detail.Text, sharesData.Detail.Trend == details.TrendNegative)
		if err == nil {
			return shares, nil
		}

		return 0, fmt.Errorf("could not parse transaction section shares to float: %w", err)
	}

	return shares, nil
}

func (b BaseModelBuilder) ExtractRateValue() (float64, error) {
	tableSections, err := b.response.SectionsTypeTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	transactionSection, err := tableSections.FindByTitle(details.SectionTitleTransaction)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	var rateData details.ResponseSectionTypeTableData

	titles := []string{
		details.TransactionDataTitleRate,
		details.TransactionDataTitleRateAlt,
		details.TransactionDataTitleRateAlt2,
	}

	for _, title := range titles {
		rateData, err = transactionSection.GetDataByTitle(title)
		if err == nil {
			break
		}
	}

	if err != nil {
		return 0, fmt.Errorf(
			"could not get transaction section rate (%s): %w", titles, err)
	}

	rate, err := ParseFloatWithComma(rateData.Detail.Text, rateData.Detail.Trend == details.TrendNegative)
	if err != nil {
		return 0, fmt.Errorf("could not parse transaction section rate to float: %w", err)
	}

	return rate, nil
}

func (b BaseModelBuilder) ExtractCommissionAmount() (float64, error) {
	tableSections, err := b.response.SectionsTypeTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	transactionSection, err := tableSections.FindByTitle(details.SectionTitleTransaction)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	commissionData, err := transactionSection.GetDataByTitle(details.TransactionDataTitleCommission)
	if err != nil {
		if !errors.Is(err, details.ErrSectionDataTitleNotFound) {
			return 0, fmt.Errorf("could not get transaction section commission: %w", err)
		}
	}

	commission, err := ParseFloatWithComma(
		commissionData.Detail.Text,
		commissionData.Detail.Trend == details.TrendNegative,
	)
	if err != nil {
		if !errors.Is(err, ErrNoMatch) {
			return 0, fmt.Errorf("could not parse transaction section commission to float: %w", err)
		}

		commission = 0
	}

	return commission, nil
}

func (b BaseModelBuilder) ExtractTotalAmount() (float64, error) {
	tableSections, err := b.response.SectionsTypeTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	transactionSection, err := tableSections.FindByTitle(details.SectionTitleTransaction)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	totalData, err := transactionSection.GetDataByTitle(details.TransactionDataTitleTotal)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section total: %w", err)
	}

	total, err := ParseFloatWithComma(totalData.Detail.Text, totalData.Detail.Trend == details.TrendNegative)
	if err != nil {
		return 0, fmt.Errorf("could not parse transaction section total to float: %w", err)
	}

	return total, nil
}

func (b BaseModelBuilder) ExtractTaxAmount() (float64, error) {
	tableSections, err := b.response.SectionsTypeTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	transactionSection, err := tableSections.FindByTitle(details.SectionTitleTransaction)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section: %w", err)
	}

	taxData, err := transactionSection.GetDataByTitle(details.TransactionDataTitleTax)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section tax amount: %w", err)
	}

	taxAmount, err := ParseFloatWithComma(taxData.Detail.Text, false)
	if err != nil {
		return 0, fmt.Errorf("could not parse transaction section tax amount to float: %w", err)
	}

	return taxAmount, nil
}

func (b BaseModelBuilder) BuildDocuments(model Model) ([]document.Model, error) {
	documents, err := b.documentsBuilder.Build(model.UUID, model.Timestamp, b.response)
	if err != nil {
		return nil, fmt.Errorf("document model builder error: %w", err)
	}

	return documents, nil
}

type PurchaseBuilder struct {
	BaseModelBuilder
}

func NewPurchaseBuilder(baseBuilder BaseModelBuilder) PurchaseBuilder {
	return PurchaseBuilder{baseBuilder}
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
	model.Documents, _ = b.BuildDocuments(model)

	return model, err
}

type SaleBuilder struct {
	PurchaseBuilder
}

func NewSaleBuilder(purchaseBuilder PurchaseBuilder) SaleBuilder {
	return SaleBuilder{purchaseBuilder}
}

func (b SaleBuilder) Build() (Model, error) {
	model, err := b.PurchaseBuilder.Build()
	if err != nil {
		return model, err
	}

	model.Type = TypeSale

	model.TaxAmount, err = b.ExtractTaxAmount()
	if err != nil {
		b.logger.WithField("id", model.UUID).Debugf("could not extract tax amount: %s", err)
	}

	model.Yield, err = b.ExtractYield()
	if err != nil {
		return model, err
	}

	model.Profit, err = b.ExtractProfitAndLoss()
	if err != nil {
		return model, err
	}

	return model, nil
}

type RoundUpBuilder struct {
	PurchaseBuilder
}

func NewRoundUpBuilder(purchaseBuilder PurchaseBuilder) RoundUpBuilder {
	return RoundUpBuilder{purchaseBuilder}
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

func NewSavebackBuilder(purchaseBuilder PurchaseBuilder) SavebackBuilder {
	return SavebackBuilder{purchaseBuilder}
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
	PurchaseBuilder
}

func NewDividendPayoutBuilder(purchaseBuilder PurchaseBuilder) DividendPayoutBuilder {
	return DividendPayoutBuilder{purchaseBuilder}
}

func (b DividendPayoutBuilder) Build() (Model, error) {
	model, err := b.PurchaseBuilder.Build()
	if err != nil {
		return model, err
	}

	model.Type = TypeDividendPayout

	return model, nil
}

type DepositBuilder struct {
	BaseModelBuilder
}

func NewDepositBuilder(baseBuilder BaseModelBuilder) DepositBuilder {
	return DepositBuilder{baseBuilder}
}

func (b DepositBuilder) Build() (Model, error) {
	var err error

	model := Model{
		UUID: b.response.ID,
		Type: TypeDeposit,
	}

	model.Status, err = b.ExtractStatus()
	if err != nil {
		return model, err
	}

	model.Timestamp, err = b.ExtractTimestamp()
	if err != nil {
		return model, err
	}

	header, err := b.response.SectionTypeHeader()
	if err != nil {
		return model, fmt.Errorf("could not get header section: %w", err)
	}

	totalAmountStr, err := ParseNumericValueFromString(header.Title)
	if err != nil {
		return model, err
	}

	model.Total, err = ParseFloatWithComma(totalAmountStr, false)
	if err != nil {
		return model, err
	}

	model.Documents, _ = b.BuildDocuments(model)

	return model, nil
}

type WithdrawBuilder struct {
	DepositBuilder
}

func NewWithdrawBuilder(depositBuilder DepositBuilder) WithdrawBuilder {
	return WithdrawBuilder{depositBuilder}
}

func (b WithdrawBuilder) Build() (Model, error) {
	model, err := b.DepositBuilder.Build()
	if err != nil {
		return model, err
	}

	model.Type = TypeWithdrawal

	return model, nil
}
