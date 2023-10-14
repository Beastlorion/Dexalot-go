package orders

import (
	"github.com/Abso1ut3Zer0/Dexalot-go/src/floats"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/price"
)

const bpsToRawFactor = 10_000.0

func Skew[T instr.Instrument](order *Limit[T], skewBps float64) {
	order.Price = price.Skew(order.Price, skewBps)
}

func RoundPricePassiveWithPrecision[T instr.Instrument](order *Limit[T], precision int) error {
	switch order.Side {
	case side.BUY:
		order.Price = floats.RoundDirectionalWithPrecision(order.Price, precision, floats.RoundDown)
		return nil
	case side.SELL:
		order.Price = floats.RoundDirectionalWithPrecision(order.Price, precision, floats.RoundUp)
		return nil
	default:
		return NoSideError
	}
}

func RoundPricePassiveWithIncrement[T instr.Instrument](order *Limit[T], increment float64) error {
	switch order.Side {
	case side.BUY:
		order.Price = floats.RoundDirectionalWithIncrement(order.Price, increment, floats.RoundDown)
		return nil
	case side.SELL:
		order.Price = floats.RoundDirectionalWithIncrement(order.Price, increment, floats.RoundUp)
		return nil
	default:
		return NoSideError
	}
}

func RoundPriceAggressiveWithPrecision[T instr.Instrument](order *Limit[T], precision int) error {
	switch order.Side {
	case side.BUY:
		order.Price = floats.RoundDirectionalWithPrecision(order.Price, precision, floats.RoundUp)
		return nil
	case side.SELL:
		order.Price = floats.RoundDirectionalWithPrecision(order.Price, precision, floats.RoundDown)
		return nil
	default:
		return NoSideError
	}
}

func RoundPriceAggressiveWithIncrement[T instr.Instrument](order *Limit[T], increment float64) error {
	switch order.Side {
	case side.BUY:
		order.Price = floats.RoundDirectionalWithIncrement(order.Price, increment, floats.RoundUp)
		return nil
	case side.SELL:
		order.Price = floats.RoundDirectionalWithIncrement(order.Price, increment, floats.RoundDown)
		return nil
	default:
		return NoSideError
	}
}

func RoundQtyWithPrecision[T instr.Instrument](order *Limit[T], precision int) {
	order.Qty = floats.RoundDirectionalWithPrecision(order.Qty, precision, floats.RoundDown)
}

func RoundQtyWithIncrement[T instr.Instrument](order *Limit[T], increment float64) {
	order.Qty = floats.RoundDirectionalWithIncrement(order.Qty, increment, floats.RoundDown)
}

func BoundQtyInDealt[T instr.Instrument](order *Limit[T], inventoryInDealt float64) error {
	switch order.Side {
	case side.BUY:
		order.Qty = floats.UpperBound(order.Qty*order.Price, inventoryInDealt) / order.Price
		return nil
	case side.SELL:
		order.Qty = floats.UpperBound(order.Qty, inventoryInDealt)
		return nil
	default:
		return NoSideError
	}
}
