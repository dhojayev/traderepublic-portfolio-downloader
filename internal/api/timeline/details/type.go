package details

type Type int

const (
	TypeUnsupported Type = iota
	TypeSaleTransaction
	TypePurchaseTransaction
	TypeDividendPayoutTransaction
	TypeRoundUpTransaction
	TypeSavebackTransaction
	TypeCardPaymentTransaction
	TypeDepositTransaction
	TypeDepositInterestReceivedTransaction
)
