package transaction

import (
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
)

var ErrUnsupportedType = errors.New("unsupported response")

type ModelBuilderFactoryInterface interface {
	Create(eventType transactions.EventType, response details.Response) (ModelBuilderInterface, error)
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

	baseBuilder := BaseModelBuilder{response: response, logger: f.logger}

	switch responseType {
	case details.TypePurchaseTransaction:
		return PurchaseBuilder{baseBuilder}, nil
	case details.TypeSaleTransaction:
		return SaleBuilder{PurchaseBuilder{baseBuilder}}, nil
	case details.TypeDividendPayoutTransaction:
		return DividendPayoutBuilder{PurchaseBuilder{baseBuilder}}, nil
	case details.TypeRoundUpTransaction:
		return RoundUpBuilder{PurchaseBuilder{baseBuilder}}, nil
	case details.TypeSavebackTransaction:
		return SavebackBuilder{PurchaseBuilder{baseBuilder}}, nil
	case details.TypeDepositTransaction:
		return DepositBuilder{baseBuilder}, nil
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
	response details.Response
	logger   *log.Logger
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

	timestamp, err := time.Parse(internal.DefaultTimeFormat, header.Data.Timestamp)
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
		b.logger.Debugf("could not get performance section yield: %s", err)
	}

	yield, err := ParseFloatWithComma(yieldData.Detail.Text, yieldData.Detail.Trend == details.TrendNegative)
	if err != nil {
		b.logger.Debugf("could not parse performance section yield to float: %s", err)
	}

	return yield, nil
}

func (b BaseModelBuilder) ExtractProfit() (float64, error) {
	tableSections, err := b.response.SectionsTypeHorizontalTable()
	if err != nil {
		return 0, fmt.Errorf("could not get table sections: %w", err)
	}

	performance, err := tableSections.FindByTitle(details.SectionTitlePerformance)
	if err != nil {
		return 0, fmt.Errorf("could not get performance section: %w", err)
	}

	profitData, err := performance.GetDataByTitle(details.PerformanceDataTitleProfit)
	if err != nil {
		b.logger.Debugf("could not get performance section profit: %s", err)
	}

	profit, err := ParseFloatWithComma(profitData.Detail.Text, profitData.Detail.Trend == details.TrendNegative)
	if err != nil {
		b.logger.Debugf("could not parse performance section profit to float: %s", err)
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

	sharesData, err := transactionSection.GetDataByTitle(details.TransactionDataTitleShares)
	if err != nil {
		b.logger.Debugf("could not get transaction section shares: %s", err)
	}

	shares, err := ParseFloatWithComma(sharesData.Detail.Text, sharesData.Detail.Trend == details.TrendNegative)
	if err != nil {
		b.logger.Debugf("could not parse transaction section shares to float: %s", err)
	}

	if strings.Contains(sharesData.Detail.Text, ".") {
		shares, err = ParseFloatWithPeriod(sharesData.Detail.Text)
		if err != nil {
			b.logger.Debugf("could not parse transaction section shares to float: %s", err)
		}
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
		rateData, err = transactionSection.GetDataByTitle(details.TransactionDataTitleRateAlt)
		if err != nil {
			return 0, fmt.Errorf(
				"could not get transaction section rate (%s): %w", titles, err)
		}
	}

	rate, err := ParseFloatWithComma(rateData.Detail.Text, rateData.Detail.Trend == details.TrendNegative)
	if err != nil {
		b.logger.Debugf("could not parse transaction section rate to float: %s", err)
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
			b.logger.Debugf("could not get transaction section commission: %s", err)
		}
	}

	commission, err := ParseFloatWithComma(
		commissionData.Detail.Text,
		commissionData.Detail.Trend == details.TrendNegative,
	)
	if err != nil {
		if !errors.Is(err, ErrNoMatch) {
			b.logger.Debugf("could not parse transaction section commission to float: %s", err)
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
		b.logger.Debugf("could not get transaction section total: %s", err)
	}

	total, err := ParseFloatWithComma(totalData.Detail.Text, totalData.Detail.Trend == details.TrendNegative)
	if err != nil {
		b.logger.Debugf("could not parse transaction section total to float: %s", err)
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
		b.logger.Debugf("could not get transaction section tax amount: %s", err)
	}

	taxAmount, err := ParseFloatWithComma(taxData.Detail.Text, false)
	if err != nil {
		b.logger.Debugf("could not parse transaction section tax amount to float: %s", err)
	}

	return taxAmount, nil
}

func (b BaseModelBuilder) BuildDocuments() ([]document.Model, error) {
	documents := make([]document.Model, 0)

	documentsSection, err := b.response.SectionTypeDocuments()
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

	model.TaxAmount, err = b.ExtractTaxAmount()
	if err != nil {
		return model, err
	}

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
	PurchaseBuilder
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

	depositAmountStr, err := ParseNumericValueFromString(header.Title)
	if err != nil {
		return model, err
	}

	model.DepositAmount, err = ParseFloatWithComma(depositAmountStr, false)
	if err != nil {
		return model, err
	}

	model.Documents, _ = b.BuildDocuments()

	return model, nil
}
