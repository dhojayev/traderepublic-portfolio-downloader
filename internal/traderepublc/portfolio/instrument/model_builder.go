package instrument

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
)

var ErrModelBuilderInsufficientDataResolved = errors.New("insufficient data resolved")

type ModelBuilderInterface interface {
	Build(response details.NormalizedResponse) (Model, error)
}

type ModelBuilder struct {
	typeResolver TypeResolverInterface
	logger       *log.Logger
}

func NewModelBuilder(
	typeResolver TypeResolverInterface,
	logger *log.Logger,
) ModelBuilder {
	return ModelBuilder{
		typeResolver: typeResolver,
		logger:       logger,
	}
}

func (b ModelBuilder) Build(response details.NormalizedResponse) (Model, error) {
	model := Model{}

	model.ISIN, _ = b.ExtractInstrumentISIN(response)
	model.Name, _ = b.ExtractInstrumentName(response)
	model.Icon, _ = b.ExtractInstrumentIcon(response)
	model.Type = b.typeResolver.Resolve(model)

	return model, nil
}

func (b ModelBuilder) ExtractInstrumentISIN(response details.NormalizedResponse) (string, error) {
	isin, valid := response.Header.Action.Payload.(string)

	if !valid || isin == "" {
		isin, _ = ExtractInstrumentISINFromIcon(response.Header.Data.Icon)
	}

	return isin, nil
}

func (b ModelBuilder) ExtractInstrumentName(response details.NormalizedResponse) (string, error) {
	asset, err := response.Overview.GetDataByTitles(
		details.OverviewDataTitleAsset,
		details.OverviewDataTitleUnderlyingAsset,
		details.OverviewDataTitleSecurity,
	)
	if err != nil {
		return "", fmt.Errorf("could not get overview section asset: %w", err)
	}

	return asset.Detail.Text, nil
}

func (b ModelBuilder) ExtractInstrumentIcon(response details.NormalizedResponse) (string, error) {
	return response.Header.Data.Icon, nil
}

func (b ModelBuilder) HandleErr(err error) error {
	if !errors.Is(err, details.ErrSectionDataTitleNotFound) {
		return err
	}

	return fmt.Errorf("%w: %w", ErrModelBuilderInsufficientDataResolved, err)
}
