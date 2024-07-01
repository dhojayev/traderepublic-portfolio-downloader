//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package instrument

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/database"
)

func ProvideModelBuilder(logger *log.Logger) ModelBuilder {
	wire.Build(
		NewModelBuilder,
		NewTypeResolver,

		wire.Bind(new(TypeResolverInterface), new(TypeResolver)),
	)

	return ModelBuilder{}
}

func ProvideInstrumentRepository(db *gorm.DB, logger *log.Logger) (*database.Repository[*Model], error) {
	return database.NewRepository[*Model](db, logger)
}
