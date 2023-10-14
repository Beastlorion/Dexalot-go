package orders

import (
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
)

func NormalizeMakerWithPrecision[T instr.Instrument](order *Limit[T], pricePrecision, qtyPrecision int) error {
	err := RoundPricePassiveWithPrecision(order, pricePrecision)
	if err != nil {
		return err
	}
	RoundQtyWithPrecision(order, qtyPrecision)
	return nil
}

func NormalizeTakerWithPrecision[T instr.Instrument](order *Limit[T], pricePrecision, qtyPrecision int) error {
	switch order.Side {
	case side.BUY:
		qtyTerm := order.Qty * order.Price // we can't breach the total term qty due to inventory constraints
		RoundPriceAggressiveWithPrecision(order, pricePrecision)
		order.Qty = qtyTerm / order.Price
		RoundQtyWithPrecision(order, qtyPrecision)
		return nil
	case side.SELL:
		RoundPriceAggressiveWithPrecision(order, pricePrecision)
		RoundQtyWithPrecision(order, qtyPrecision)
		return nil
	default:
		return NoSideError
	}
}

func NormalizeMakerWithIncrement[T instr.Instrument](order *Limit[T], priceIncrement, qtyIncrement float64) error {
	err := RoundPricePassiveWithIncrement(order, priceIncrement)
	if err != nil {
		return err
	}
	RoundQtyWithIncrement(order, qtyIncrement)
	return nil
}

func NormalizeTakerWithIncrement[T instr.Instrument](order *Limit[T], priceIncrement, qtyIncrement float64) error {
	switch order.Side {
	case side.BUY:
		qtyTerm := order.Qty * order.Price // we can't breach the total term qty due to inventory constraints
		RoundPriceAggressiveWithIncrement(order, priceIncrement)
		order.Qty = qtyTerm / order.Price
		RoundQtyWithIncrement(order, qtyIncrement)
		return nil
	case side.SELL:
		RoundPriceAggressiveWithIncrement(order, priceIncrement)
		RoundQtyWithIncrement(order, qtyIncrement)
		return nil
	default:
		return NoSideError
	}
}

func NormalizeMakerWithRefData[T instr.Instrument](order *Limit[T], refData *refdata.Composite) error {
	if refData == nil || refData.PriceQty == nil {
		return refdata.NilRefDataError
	}
	switch refData.PriceQty.(type) {
		case *refdata.Precision:
			rd, _ := refData.PriceQty.(*refdata.Precision)
			return NormalizeMakerWithPrecision(order, rd.PricePrecision, rd.QtyPrecision)
		case *refdata.Increment:
			rd, _ := refData.PriceQty.(*refdata.Increment)
			return NormalizeMakerWithIncrement(order, rd.PriceIncrement, rd.QtyIncrement)
		default:
			return refdata.UnsupportedRefDataError
	}
}

func NormalizeTakerWithRefData[T instr.Instrument](order *Limit[T], refData *refdata.Composite) error {
	if refData == nil || refData.PriceQty == nil {
		return refdata.NilRefDataError
	}
	switch refData.PriceQty.(type) {
	case *refdata.Precision:
		rd, _ := refData.PriceQty.(*refdata.Precision)
		return NormalizeTakerWithPrecision(order, rd.PricePrecision, rd.QtyPrecision)
	case *refdata.Increment:
		rd, _ := refData.PriceQty.(*refdata.Increment)
		return NormalizeTakerWithIncrement(order, rd.PriceIncrement, rd.QtyIncrement)
	default:
		return refdata.UnsupportedRefDataError
	}
}
