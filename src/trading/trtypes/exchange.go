package trtypes

type Exchange string

const (
	NilExchange Exchange = ""
	Dexalot    Exchange = "DEXALOT"
)

type ExchangeMetaData struct {
	Exchange Exchange
	OrderID  *string // Exchange order ID, optional until we receive it from the exchange
}
