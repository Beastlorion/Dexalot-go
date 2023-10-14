package exchange

type Exchange string

const (
	NilExchange Exchange = ""
	Dexalot     Exchange = "DEXALOT"
	Binance     Exchange = "BINANCE"
	Coinbase    Exchange = "COINBASE"
)

type MetaData struct {
	Exchange Exchange
	OrderID  *string // Exchange order ID, optional until we receive it from the exchange
}
