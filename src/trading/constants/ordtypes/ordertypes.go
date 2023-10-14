package ordtypes

type OrderType int

const (
	NilOrderType = OrderType(iota)
	Market
	Limit
	Stop
	StopLimit
)
