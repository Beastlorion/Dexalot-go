package refdata

import (
	"errors"
	"time"
)

var ReferenceDataNotFoundError = errors.New("reference data not found")
var NoTradeSizeLimitTypeError = errors.New("no trade size limit type")
var MinimumSizeError = errors.New("order size is less than minimum trade size")
var MaximumSizeError = errors.New("order size is greater than maximum trade size")
var NilRefDataError = errors.New("nil reference data")
var UnsupportedRefDataError = errors.New("unsupported reference data")

// We have to use type assertions or type selecting for mutually exclusive types
// This is just a marker to help us know we expect one of the Interface types defined below.
type Interface interface{}

type Precision struct {
	PricePrecision int
	QtyPrecision   int
}

type Increment struct {
	PriceIncrement float64
	QtyIncrement   float64
}

type BaseTradeSizeLimit struct {
	BaseMinTradeSize float64
	BaseMaxTradeSize float64
}

type TermTradeSizeLimit struct {
	TermMinTradeSize float64
	TermMaxTradeSize float64
}

type Composite struct {
	PriceQty          Interface
	TradeSizeLimit    Interface
	MinQuoteTime      time.Duration
	ReplenishmentRate time.Duration
	UsePostOnly       bool
}
