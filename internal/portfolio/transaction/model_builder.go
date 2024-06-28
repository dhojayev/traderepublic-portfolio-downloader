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

var (
	ErrModelBuilderUnsupportedType          = errors.New("unsupported response")
	ErrModelBuilderInsufficientDataResolved = errors.New("insufficient data resolved")
	ErrModelBuilderUnknownType              = errors.New("unknown response")
)

type ModelBuilderFactoryInterface interface {
	Create(eventType transactions.EventType, response details.NormalizedResponse) (ModelBuilderInterface, error)
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
	response details.NormalizedResponse,
) (ModelBuilderInterface, error) {
	responseType, err := f.resolver.Resolve(eventType, response)
	if err != nil {
		if errors.Is(err, details.ErrTypeResolverUnsupportedType) {
			return nil, ErrModelBuilderUnsupportedType
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
		return nil, ErrModelBuilderUnsupportedType
	}

	return nil, ErrModelBuilderUnknownType
}

type ModelBuilderInterface interface {
	Build() (Model, error)
}

type BaseModelBuilder struct {
	response         details.NormalizedResponse
	documentsBuilder document.ModelBuilderInterface
	logger           *log.Logger
}

func NewBaseModelBuilder(
	response details.NormalizedResponse,
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
	return b.response.Header.Data.Status, nil
}

func (b BaseModelBuilder) ExtractInstrumentIcon() (string, error) {
	return b.response.Header.Data.Icon, nil
}

func (b BaseModelBuilder) ExtractTimestamp() (time.Time, error) {
	timestamp, err := time.Parse(details.ResponseTimeFormat, b.response.Header.Data.Timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse header section timestamp: %w", err)
	}

	return timestamp, nil
}

func (b BaseModelBuilder) ExtractInstrumentISIN() (string, error) {
	isin, valid := b.response.Header.Action.Payload.(string)

	if !valid || isin == "" {
		isin, _ = ExtractInstrumentNameFromIcon(b.response.Header.Data.Icon)
	}

	return isin, nil
}

func (b BaseModelBuilder) ExtractInstrumentName() (string, error) {
	asset, err := b.response.Overview.GetDataByTitles(
		details.OverviewDataTitleAsset,
		details.OverviewDataTitleUnderlyingAsset,
		details.OverviewDataTitleSecurity,
	)
	if err != nil {
		return "", fmt.Errorf("could not get overview section asset: %w", err)
	}

	return asset.Detail.Text, nil
}

func (b BaseModelBuilder) ExtractSharesAmount() (float64, error) {
	sharesData, err := b.response.Transaction.GetDataByTitles(
		details.TransactionDataTitleShares,
		details.TransactionDataTitleSharesAlt,
	)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction section shares: %w", err)
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
	rateData, err := b.response.Transaction.GetDataByTitles(
		details.TransactionDataTitleRate,
		details.TransactionDataTitleRateAlt,
		details.TransactionDataTitleRateAlt2,
		details.TransactionDataTitleRateAlt3,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"could not get transaction section rate: %w", err)
	}

	rate, err := ParseFloatWithComma(rateData.Detail.Text, rateData.Detail.Trend == details.TrendNegative)
	if err != nil {
		return 0, fmt.Errorf("could not parse transaction section rate to float: %w", err)
	}

	return rate, nil
}

func (b BaseModelBuilder) ExtractCommissionAmount() (float64, error) {
	commissionData, err := b.response.Transaction.GetDataByTitles(details.TransactionDataTitleCommission)
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
	totalData, err := b.response.Transaction.GetDataByTitles(details.TransactionDataTitleTotal)
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
	taxData, err := b.response.Transaction.GetDataByTitles(details.TransactionDataTitleTax)
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

func (b BaseModelBuilder) HandleErr(err error) error {
	if !errors.Is(err, details.ErrSectionDataTitleNotFound) {
		return err
	}

	return fmt.Errorf("%w: %w", ErrModelBuilderInsufficientDataResolved, err)
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
		return model, b.HandleErr(err)
	}

	model.Timestamp, err = b.ExtractTimestamp()
	if err != nil {
		return model, b.HandleErr(err)
	}

	model.Instrument.ISIN, err = b.ExtractInstrumentISIN()
	if err != nil {
		return model, b.HandleErr(err)
	}

	model.Instrument.Name, err = b.ExtractInstrumentName()
	if err != nil {
		return model, b.HandleErr(err)
	}

	model.Shares, err = b.ExtractSharesAmount()
	if err != nil {
		return model, b.HandleErr(err)
	}

	model.Rate, err = b.ExtractRateValue()
	if err != nil {
		return model, b.HandleErr(err)
	}

	model.Commission, err = b.ExtractCommissionAmount()
	if err != nil {
		return model, b.HandleErr(err)
	}

	model.Total, err = b.ExtractTotalAmount()
	if err != nil {
		return model, b.HandleErr(err)
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

func (b SaleBuilder) ExtractPerformanceFloatVal(titles ...string) (float64, error) {
	stringData, err := b.response.Performance.GetDataByTitles(titles...)
	if err != nil {
		return 0, fmt.Errorf("could not get performance section data by titles %v: %w", titles, err)
	}

	floatData, err := ParseFloatWithComma(stringData.Detail.Text, stringData.Detail.Trend == details.TrendNegative)
	if err != nil {
		return 0, fmt.Errorf("could not parse performance section data to float by titles %v: %w", titles, err)
	}

	return floatData, nil
}

func (b SaleBuilder) ExtractYield() (float64, error) {
	return b.ExtractPerformanceFloatVal(details.PerformanceDataTitleYield)
}

func (b SaleBuilder) ExtractProfitAndLoss() (float64, error) {
	return b.ExtractPerformanceFloatVal(details.PerformanceDataTitleProfit, details.PerformanceDataTitleLoss)
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
		return model, b.HandleErr(err)
	}

	model.Profit, err = b.ExtractProfitAndLoss()
	if err != nil {
		return model, b.HandleErr(err)
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
		return model, b.HandleErr(err)
	}

	model.Timestamp, err = b.ExtractTimestamp()
	if err != nil {
		return model, err
	}

	model.Total, err = b.ExtractTotalAmount()
	if err != nil {
		return model, b.HandleErr(err)
	}

	model.Documents, _ = b.BuildDocuments(model)

	return model, nil
}

func (b DepositBuilder) ExtractTotalAmount() (float64, error) {
	totalAmountStr, err := ParseNumericValueFromString(b.response.Header.Title)
	if err != nil {
		return 0, err
	}

	total, err := ParseFloatWithComma(totalAmountStr, false)
	if err != nil {
		return total, err
	}

	return total, nil
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
