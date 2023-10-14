package orders

import (
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
)

func ValidateLimit[T instr.Instrument](order *Limit[T]) error {
	if order == nil {
		return NoOrderError
	} else if order.ExchangeMetaData == nil {
		return NoExchangeError
	} else if order.ExchangeMetaData.Exchange == exchange.NilExchange {
		return NoExchangeError
	} else if order.Side == side.NilSide {
		return NoSideError
	} else if order.TimeInForce == tif.NilTimeInForce {
		return NoTimeInForceError
	} else if !order.Instrument.NegativePriceAllowed() && order.Price <= 0 {
		return NonPositivePriceError
	} else if order.Qty <= 0 {
		return NonPositiveQtyError
	}
	return nil
}

func ValidateLimitWithBaseConstraints[T instr.Instrument](order *Limit[T], baseMinTradeSize, baseMaxTradeSize float64) error {
	if err := ValidateLimit(order); err != nil {
		return err
	} else if order.Qty < baseMinTradeSize {
		return refdata.MinimumSizeError
	} else if order.Qty > baseMaxTradeSize {
		return refdata.MaximumSizeError
	}
	return nil
}

func ValidateLimitWithTermConstraints[T instr.Instrument](order *Limit[T], termMinTradeSize, termMaxTradeSize float64) error {
	if err := ValidateLimit(order); err != nil {
		return err
	} else if order.Qty*order.Price < termMinTradeSize {
		return refdata.MinimumSizeError
	} else if order.Qty*order.Price > termMaxTradeSize {
		return refdata.MaximumSizeError
	}
	return nil
}

func ValidateLimitWithRefData[T instr.Instrument](order *Limit[T], refData *refdata.Composite) error {
	if refData == nil || refData.TradeSizeLimit == nil {
		return refdata.NilRefDataError
	}
	switch refData.TradeSizeLimit.(type) {
	case *refdata.BaseTradeSizeLimit:
		rd, _ := refData.TradeSizeLimit.(*refdata.BaseTradeSizeLimit)
		return ValidateLimitWithBaseConstraints(order, rd.BaseMinTradeSize, rd.BaseMaxTradeSize)
	case *refdata.TermTradeSizeLimit:
		rd, _ := refData.TradeSizeLimit.(*refdata.TermTradeSizeLimit)
		return ValidateLimitWithTermConstraints(order, rd.TermMinTradeSize, rd.TermMaxTradeSize)
	default:
		return refdata.UnsupportedRefDataError
	}
}
