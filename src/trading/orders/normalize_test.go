package orders_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
)

func TestNormalizeMakerPrecision(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Side:  side.SELL,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeMakerWithPrecision(order, 2, 3)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 2, 2)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 2, 1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 2, 0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 1, 0)
	if order.Price != 100.2 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.2", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 0, 0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPricePassive() = %f; want 101.0", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeMakerWithPrecision(order, 2, 3)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 2, 2)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 2, 1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 2, 0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 1, 0)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.1", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithPrecision(order, 0, 0)
	if order.Price != 100.0 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.0", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.13,
		Qty:   100.134,
	}

	if err := orders.NormalizeMakerWithPrecision(order, 2, 3); err == nil {
		t.Errorf("RoundOrderPricePassive() = %v; want nil", err)
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}
}

func TestNormalizeTakerPrecision(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Side:  side.SELL,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeTakerWithPrecision(order, 2, 3)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 2, 2)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 2, 1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 2, 0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 1, 0)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.1", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 0, 0)
	if order.Price != 100.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.0", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeTakerWithPrecision(order, 2, 3)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 2, 2)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 2, 1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 2, 0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 1, 0)
	if order.Price != 100.2 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.1", order.Price)
	} else if order.Qty != 99.0 {
		t.Errorf("RoundOrderQty() = %f; want 99.0", order.Qty)
	}

	orders.NormalizeTakerWithPrecision(order, 0, 0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.0", order.Price)
	} else if order.Qty != 98.0 {
		t.Errorf("RoundOrderQty() = %f; want 98.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.13,
		Qty:   100.134,
	}

	if err := orders.NormalizeTakerWithPrecision(order, 2, 3); err == nil {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type: %v", err)
	}
}

func TestNormalizeMakerIncrement(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Side:  side.SELL,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.001)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.01)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.1, 0.0)
	if order.Price != 100.2 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.2", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.0, 0.0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPricePassive() = %f; want 101.0", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.001)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.01)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.01, 0.0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.1, 0.0)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.1", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeMakerWithIncrement(order, 0.0, 0.0)
	if order.Price != 100.0 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.0", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.13,
		Qty:   100.134,
	}

	if err := orders.NormalizeMakerWithIncrement(order, 0.01, 0.001); err == nil {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}
}

func TestNormalizeTakerIncrement(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Side:  side.SELL,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.001)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.01)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.1, 0.0)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.1", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.0, 0.0)
	if order.Price != 100.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.0", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.13,
		Qty:   100.134,
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.001)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.134 {
		t.Errorf("RoundOrderQty() = %f; want 100.134", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.01)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.1)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.01, 0.0)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	} else if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.1, 0.0)
	if order.Price != 100.2 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.1", order.Price)
	} else if order.Qty != 99.0 {
		t.Errorf("RoundOrderQty() = %f; want 99.0", order.Qty)
	}

	orders.NormalizeTakerWithIncrement(order, 0.0, 0.0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.0", order.Price)
	} else if order.Qty != 98.0 {
		t.Errorf("RoundOrderQty() = %f; want 98.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.0,
		Qty:   100.0,
	}

	if err := orders.NormalizeTakerWithIncrement(order, 0.01, 0.001); err == nil {
		t.Errorf("Expected error")
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}
}

func TestNormalizeMakerWithRefData(t *testing.T) {
	order := &orders.Limit[instr.Spot]{}
	refData := &refdata.Composite{}

	// only check the error type
	if err := orders.NormalizeMakerWithRefData(order, refData); err == nil {
		t.Errorf("Expected error")
	} else if err != refdata.NilRefDataError {
		t.Errorf("Wrong error type")
	}

	refData = &refdata.Composite{PriceQty: struct{}{}}
	if err := orders.NormalizeMakerWithRefData(order, refData); err == nil {
	    t.Errorf("Expected error")
	} else if err != refdata.UnsupportedRefDataError {
	    t.Errorf("Wrong error type")
	}
}

func TestNormalizeTakerWithRefData(t *testing.T) {
	order := &orders.Limit[instr.Spot]{}
	refData := &refdata.Composite{}

	// only check the error type
	if err := orders.NormalizeTakerWithRefData(order, refData); err == nil {
		t.Errorf("Expected error")
	} else if err != refdata.NilRefDataError {
		t.Errorf("Wrong error type")
	}

	refData = &refdata.Composite{PriceQty: struct{}{}}
	// only check the error type
	if err := orders.NormalizeTakerWithRefData(order, refData); err == nil {
	    t.Errorf("Expected error")
	} else if err != refdata.UnsupportedRefDataError {
	    t.Errorf("Wrong error type")
	}
}
