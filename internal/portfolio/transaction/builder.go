package transaction

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
)

var ErrUnsupportedResponse = errors.New("unsupported response")

type BuilderInterface interface {
	FromResponse(response details.Response) (Model, error)
}

type Builder struct {
	resolver TypeResolverInterface
	logger   *log.Logger
}

func NewBuilder(resolver TypeResolverInterface, logger *log.Logger) Builder {
	return Builder{
		resolver: resolver,
		logger:   logger,
	}
}

func (b Builder) FromResponse(response details.Response) (Model, error) {
	resolvedType, err := b.resolver.Resolve(response)
	if err != nil {
		return Model{}, fmt.Errorf("resolver error: %w", err)
	}

	switch resolvedType {
	case TypePurchaseTransaction:
		return b.Build(TypePurchase, response)
	case TypeSaleTransaction:
		return b.Build(TypeSale, response)
	case TypeDividendPayoutTransaction:
		return b.Build(TypeDividendPayout, response)
	case TypeRoundUpTransaction:
		return b.Build(TypeRoundUp, response)
	case TypeSavebackTransaction:
		return b.Build(TypeSaveback, response)
	case TypeUnsupported, TypeCardPaymentTransaction:
	}

	return Model{}, ErrUnsupportedResponse
}

func (b Builder) Build(transactionType string, response details.Response) (Model, error) {
	var err error

	transaction := Model{
		UUID: response.ID,
		Type: transactionType,
	}

	transaction.Status, transaction.Instrument.ISIN, transaction.Timestamp, err = b.GetHeaderData(response)
	if err != nil {
		return transaction, err
	}

	transaction.Instrument.Name, err = b.GetOverviewData(response)
	if err != nil {
		return transaction, err
	}

	transaction.Yield, transaction.Profit, _ = b.GetPerformanceData(response)

	transaction.Shares, transaction.Rate, transaction.Commission, transaction.Total, err = b.GetTransactionData(response)
	if err != nil {
		return transaction, err
	}

	documents, err := b.BuildDocuments(response)
	if err != nil {
		return transaction, err
	}

	transaction.Documents = documents

	return transaction, nil
}

// GetHeaderData Returns Status, ISIN, Timestamp and error.
func (b Builder) GetHeaderData(response details.Response) (string, string, time.Time, error) {
	var status, isin string

	var timestamp time.Time

	header, err := response.HeaderSection()
	if err != nil {
		return status, isin, timestamp, fmt.Errorf("could not get details header %w", err)
	}

	status = header.Data.Status

	timestamp, err = time.Parse("2006-01-02T15:04:05-0700", header.Data.Timestamp)
	if err != nil {
		b.logger.Debugf("could not parse details timestamp: %s", err)
	}

	isin, _ = header.Action.Payload.(string)
	if isin == "" {
		isin, _ = ExtractInstrumentNameFromIcon(header.Data.Icon)
	}

	return status, isin, timestamp, nil
}

// GetOverviewData Returns Instrument name and error.
func (b Builder) GetOverviewData(response details.Response) (string, error) {
	var instrumentName string

	overview, err := response.OverviewSection()
	if err != nil {
		return instrumentName, fmt.Errorf("error getting overview: %w", err)
	}

	asset, err := overview.Asset()
	if err != nil {
		return instrumentName, fmt.Errorf("error getting overview asset: %w", err)
	}

	instrumentName = asset.Detail.Text

	return instrumentName, nil
}

// GetPerformanceData Returns Yield, Profit and error.
func (b Builder) GetPerformanceData(response details.Response) (float64, float64, error) {
	var yield, profit float64

	performance, err := response.PerformanceSection()
	if err != nil {
		return yield, profit, fmt.Errorf("could get performance section: %w", err)
	}

	yieldData, err := performance.Yield()
	if err != nil {
		b.logger.Debugf("could get yield: %s", err)
	}

	yield, err = ParseFloatWithComma(yieldData.Detail.Text, yieldData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse float value from yield: %s", err)
	}

	profitData, err := performance.Profit()
	if err != nil {
		b.logger.Debugf("could get profit: %s", err)
	}

	profit, err = ParseFloatWithComma(profitData.Detail.Text, profitData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse float value from profit: %s", err)
	}

	return yield, profit, nil
}

// GetTransactionData Returns Shares, Rate, Commission, Total and error.
//
//nolint:cyclop,funlen
func (b Builder) GetTransactionData(response details.Response) (float64, float64, float64, float64, error) {
	var shares, rate, commission, total float64

	transactionSection, err := response.TransactionSection()
	if err != nil {
		return shares, rate, commission, total, fmt.Errorf("could not get transaction section: %w", err)
	}

	sharesData, err := transactionSection.Shares()
	if err != nil {
		b.logger.Debugf("could not get shares: %s", err)
	}

	shares, err = ParseFloatWithComma(sharesData.Detail.Text, sharesData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse float value from shares: %s", err)
	}

	if sharesData.HasSharesWithPeriod() {
		shares, err = ParseFloatWithPeriod(sharesData.Detail.Text)
		if err != nil {
			b.logger.Debugf("could not parse float value from shares: %s", err)
		}
	}

	rateData, err := transactionSection.Rate()
	if err != nil {
		b.logger.Debugf("could not get rate: %s", err)
	}

	rate, err = ParseFloatWithComma(rateData.Detail.Text, rateData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse float value from rate: %s", err)
	}

	commissionData, err := transactionSection.Commission()
	if err != nil {
		if !errors.Is(err, details.ErrSectionDataEntryNotFound) {
			b.logger.Debugf("could not get commission: %s", err)
		}
	}

	commission, err = ParseFloatWithComma(commissionData.Detail.Text, commissionData.Detail.IsTrendNegative())
	if err != nil {
		if !errors.Is(err, ErrNoMatch) {
			b.logger.Debugf("could not parse float value from commission: %s", err)
		}

		commission = 0
	}

	totalData, err := transactionSection.Total()
	if err != nil {
		b.logger.Debugf("could not get total: %s", err)
	}

	total, err = ParseFloatWithComma(totalData.Detail.Text, totalData.Detail.IsTrendNegative())
	if err != nil {
		b.logger.Debugf("could not parse float value from total: %s", err)
	}

	return shares, rate, commission, total, nil
}

func (b Builder) BuildDocuments(response details.Response) ([]Document, error) {
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

		documents = append(documents, NewDocument(document.ID, url, document.Detail, document.Title))
	}

	return documents, nil
}
