package instrument

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	TypeStocks         Type = "Stocks"
	TypeETF            Type = "ETF"
	TypeCryptocurrency Type = "Cryptocurrency"
	TypeLending        Type = "Lending"
	TypeCash           Type = "Cash"
	TypeOther          Type = "Other"
)

type Type string

type TypeResolverInterface interface {
	Resolve(instrument Model) Type
}

type TypeResolver struct {
	logger *log.Logger
}

func NewTypeResolver(logger *log.Logger) TypeResolver {
	return TypeResolver{
		logger: logger,
	}
}

func (r TypeResolver) Resolve(instrument Model) Type {
	switch {
	case instrument.Name == "", strings.HasPrefix(instrument.Name, isinPrefixCash):
		return TypeCash
	case strings.HasSuffix(instrument.Name, isinSuffixDist), strings.HasSuffix(instrument.Name, isinSuffixAcc):
		return TypeETF
	case strings.HasPrefix(instrument.ISIN, isinPrefixCrypto):
		return TypeCryptocurrency
	case strings.HasPrefix(instrument.ISIN, isinPrefixLending):
		return TypeLending
	}

	return TypeOther
}
