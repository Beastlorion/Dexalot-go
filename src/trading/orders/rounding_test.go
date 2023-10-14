package orders_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
)

func TestSkew(t *testing.T) {
	order := &orders.Limit[instr.Spot]{Price: 100.0}

	orders.Skew(order, 100.0)
	if order.Price != 101.0 {
		t.Errorf("SkewOrderPrice() = %f; want 101.0", order.Price)
	}

	order = &orders.Limit[instr.Spot]{Price: 100.0}
	orders.Skew(order, -100.0)
	if order.Price != 99.0 {
		t.Errorf("SkewOrderPrice() = %f; want 99.0", order.Price)
	}

	order = &orders.Limit[instr.Spot]{Price: 100.0}
	orders.Skew(order, 0.0)
	if order.Price != 100.0 {
		t.Errorf("SkewOrderPrice() = %f; want 100.0", order.Price)
	}

	order = &orders.Limit[instr.Spot]{Price: 100.0}
	orders.Skew(order, 500.0)
	if order.Price != 105.0 {
		t.Errorf("SkewOrderPrice() = %f; want 105.0", order.Price)
	}
}

func TestRoundPricePassiveWithPrecision(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.13,
	}

	orders.RoundPricePassiveWithPrecision(order, 2)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 1)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.1", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 2)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.1", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.154,
		Side:  side.BUY,
	}
	orders.RoundPricePassiveWithPrecision(order, 4)
	if order.Price != 100.154 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.154", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 3)
	if order.Price != 100.154 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.154", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 2)
	if order.Price != 100.15 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.15", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
		Side:  side.SELL,
	}

	orders.RoundPricePassiveWithPrecision(order, 4)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 3)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 2)
	if order.Price != 100.50 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.50", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 1)
	if order.Price != 100.5 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.5", order.Price)
	}

	orders.RoundPricePassiveWithPrecision(order, 0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPricePassive() = %f; want 101.0", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
	}

	err := orders.RoundPricePassiveWithPrecision(order, 4)
	if err == nil {
		t.Errorf("RoundOrderPricePassive() = %f; want error", order.Price)
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}
}

func TestRoundPricePassiveIncrement(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.13,
	}

	orders.RoundPricePassiveWithIncrement(order, 0.01)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.13", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.1)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.1", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.01)
	if order.Price != 100.1 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.1", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.154,
		Side:  side.BUY,
	}

	orders.RoundPricePassiveWithIncrement(order, 0.0001)
	if order.Price != 100.154 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.154", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.001)
	if order.Price != 100.154 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.154", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.01)
	if order.Price != 100.15 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.15", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
		Side:  side.SELL,
	}
	orders.RoundPricePassiveWithIncrement(order, 0.0001)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.001)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.01)
	if order.Price != 100.50 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.50", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.1)
	if order.Price != 100.5 {
		t.Errorf("RoundOrderPricePassive() = %f; want 100.5", order.Price)
	}

	orders.RoundPricePassiveWithIncrement(order, 0.0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPricePassive() = %f; want 101.0", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
	}

	err := orders.RoundPricePassiveWithIncrement(order, 0.001)
	if err == nil {
		t.Errorf("RoundOrderPricePassive() = %f; want error", order.Price)
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}
}

func TestRoundPriceAggressive(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Price: 100.13,
		Side:  side.BUY,
	}

	orders.RoundPriceAggressiveWithPrecision(order, 3)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	}

	orders.RoundPriceAggressiveWithPrecision(order, 2)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	}

	orders.RoundPriceAggressiveWithPrecision(order, 1)
	if order.Price != 100.2 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.2", order.Price)
	}

	orders.RoundPriceAggressiveWithPrecision(order, 0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 101.0", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
		Side:  side.SELL,
	}
	orders.RoundPriceAggressiveWithPrecision(order, 4)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPriceAggressiveWithPrecision(order, 3)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPriceAggressiveWithPrecision(order, 2)
	if order.Price != 100.49 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.49", order.Price)
	}

	orders.RoundPriceAggressiveWithPrecision(order, 1)
	if order.Price != 100.4 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
		Side:  side.SELL,
	}

	orders.RoundPriceAggressiveWithPrecision(order, 0)
	if order.Price != 100.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
	}

	err := orders.RoundPriceAggressiveWithPrecision(order, 3)
	if err == nil {
		t.Errorf("RoundOrderPricePassive() = %f; want error", order.Price)
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}
}

func TestRoundPriceAggressiveIncrement(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Price: 100.13,
		Side:  side.BUY,
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.001)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.01)
	if order.Price != 100.13 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.13", order.Price)
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.1)
	if order.Price != 100.2 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.2", order.Price)
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.0)
	if order.Price != 101.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 101.0", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
		Side:  side.SELL,
	}
	orders.RoundPriceAggressiveWithIncrement(order, 0.0001)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.001)
	if order.Price != 100.4980 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4980", order.Price)
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.01)
	if order.Price != 100.49 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.49", order.Price)
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.1)
	if order.Price != 100.4 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4", order.Price)
	}

	orders.RoundPriceAggressiveWithIncrement(order, 0.0)
	if order.Price != 100.0 {
		t.Errorf("RoundOrderPriceAggressive() = %f; want 100.4", order.Price)
	}

	order = &orders.Limit[instr.Spot]{
		Price: 100.4980,
	}
	if err := orders.RoundPriceAggressiveWithIncrement(order, 0.001); err == nil {
		t.Errorf("RoundOrderPriceAggressive() = %f; want error", order.Price)
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}
}

func TestBoundOrderQtyInDealt(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Side: side.SELL,
		Qty:  100.1357,
	}
	orders.BoundQtyInDealt(order, 100.0)

	if order.Qty != 100.0 {
		t.Errorf("BoundOrderQty() = %f; want 100.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Side: side.SELL,
		Qty:  100.1357,
	}
	orders.BoundQtyInDealt(order, 100.15)

	if order.Qty != 100.1357 {
		t.Errorf("BoundOrderQty() = %f; want 100.1357", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.0,
		Qty:   1.5,
	}

	orders.BoundQtyInDealt(order, 100.0)
	if order.Qty != 1.0 {
		t.Errorf("BoundOrderQty() = %f; want 1.0", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Side:  side.BUY,
		Price: 100.0,
		Qty:   1.5,
	}

	orders.BoundQtyInDealt(order, 200.0)
	if order.Qty != 1.5 {
		t.Errorf("BoundOrderQty() = %f; want 1.5", order.Qty)
	}

	order = &orders.Limit[instr.Spot]{
		Qty: 100.1357,
	}

	if err := orders.BoundQtyInDealt(order, 200.0); err == nil {
		t.Errorf("BoundOrderQty() = %f; want error", order.Qty)
	} else if err != orders.NoSideError {
		t.Errorf("Wronge error type")
	}
}

func TestRoundQtyWithPrecision(t *testing.T) {
	order := &orders.Limit[instr.Spot]{Qty: 100.13}

	orders.RoundQtyWithPrecision(order, 3)
	if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.RoundQtyWithPrecision(order, 2)
	if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.RoundQtyWithPrecision(order, 1)
	if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.RoundQtyWithPrecision(order, 0)
	if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}
}

func TestRoundQtyWithIncrement(t *testing.T) {
	order := &orders.Limit[instr.Spot]{Qty: 100.13}

	orders.RoundQtyWithIncrement(order, 0.001)
	if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.RoundQtyWithIncrement(order, 0.01)
	if order.Qty != 100.13 {
		t.Errorf("RoundOrderQty() = %f; want 100.13", order.Qty)
	}

	orders.RoundQtyWithIncrement(order, 0.1)
	if order.Qty != 100.1 {
		t.Errorf("RoundOrderQty() = %f; want 100.1", order.Qty)
	}

	orders.RoundQtyWithIncrement(order, 0.0)
	if order.Qty != 100.0 {
		t.Errorf("RoundOrderQty() = %f; want 100.0", order.Qty)
	}
}
