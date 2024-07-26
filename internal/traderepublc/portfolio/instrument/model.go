package instrument

import (
	"fmt"
)

const (
	isinPrefixLending = "XS"
	isinPrefixCrypto  = "XF000"
	isinSuffixDist    = "(Dist)"
	isinSuffixAcc     = "(Acc)"
)

type Model struct {
	ISIN string `gorm:"primaryKey"`
	Name string
	Icon string
	Type Type
}

func NewModel(isin, name, icon string, instrumentType Type) Model {
	return Model{
		ISIN: isin,
		Name: name,
		Type: instrumentType,
		Icon: icon,
	}
}

func (i Model) IconURL() string {
	return fmt.Sprintf("https://assets.traderepublic.com/img/%s/light.min.svg", i.Icon)
}

func (Model) TableName() string {
	return "instruments"
}
