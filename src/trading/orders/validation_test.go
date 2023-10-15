package orders_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
)

func TestValidateLimitOrder(t *testing.T) {
	if err := orders.ValidateLimit[instr.Spot](nil); err == nil {
		t.Errorf("Expected error")
	} else if err != orders.NoOrderError {
		t.Errorf("Wrong error type")
	}

	order := &orders.Limit[instr.Spot]{
		Instrument: instr.Spot{Base: instr.BTC, Term: instr.USD},
	}

	if err := orders.ValidateLimit(order); err == nil {
		t.Errorf("Exected error")
	} else if err != orders.NoExchangeError {
		t.Errorf("Wrong error type")
	}

	order.ExchangeMetaData = &exchange.MetaData{}

	if err := orders.ValidateLimit(order); err == nil {
		t.Errorf("Expected error")
	} else if err != orders.NoExchangeError {
		t.Errorf("Wrong error type")
	}

	order.ExchangeMetaData.Exchange = exchange.Coinbase

	if err := orders.ValidateLimit(order); err == nil {
		t.Errorf("Expected error")
	} else if err != orders.NoSideError {
		t.Errorf("Wrong error type")
	}

	order.Side = side.BUY

	if err := orders.ValidateLimit(order); err == nil {
		t.Errorf("Expected error")
	} else if err != orders.NoTimeInForceError {
		t.Errorf("Wrong error type")
	}

	order.TimeInForce = tif.GTC

	if err := orders.ValidateLimit(order); err == nil {
		t.Errorf("Expected error")
	} else if err != orders.NonPositivePriceError {
		t.Errorf("Wrong error type")
	}

	order.Price = 1000.0

	if err := orders.ValidateLimit(order); err == nil {
		t.Errorf("Expected error")
	} else if err != orders.NonPositiveQtyError {
		t.Errorf("Wrong error type")
	}

	order.Qty = 100.0

	if err := orders.ValidateLimit(order); err != nil {
		t.Errorf("Expected no error")
	}
}

func TestValidateLimitOrderWithBaseConstraints(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.BTC, Term: instr.USD},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		Side:             side.BUY,
		Price:            100.0,
		Qty:              100.0,
		TimeInForce:      tif.GTC,
	}

	err := orders.ValidateLimitWithBaseConstraints(order, 0.1, 1000.0)
	if err != nil {
		t.Errorf("ValidateOrderWithBaseConstraints() = %v; want nil", err)
	}

	err = orders.ValidateLimitWithBaseConstraints(order, 0.1, 100.0)
	if err != nil {
		t.Errorf("ValidateOrderWithBaseConstraints() = %v; want nil", err)
	}

	err = orders.ValidateLimitWithBaseConstraints(order, 0.1, 99.0)
	if err == nil {
		t.Errorf("ValidateOrderWithBaseConstraints() = nil; want error")
	}

	order.Qty = 0.01
	err = orders.ValidateLimitWithBaseConstraints(order, 0.1, 1000.0)
	if err == nil {
		t.Errorf("ValidateOrderWithBaseConstraints() = nil; want error")
	}

	order.Qty = 100.0
	order.Price = 0.0
	err = orders.ValidateLimitWithBaseConstraints(order, 0.1, 1000.0)
	if err == nil {
		t.Errorf("ValidateOrderWithBaseConstraints() = nil; want error")
	}

	order.Price = -1.0
	err = orders.ValidateLimitWithBaseConstraints(order, 0.1, 1000.0)
	if err == nil {
		t.Errorf("ValidateOrderWithBaseConstraints() = nil; want error")
	}
}

func TestValidateLimitOrderWithTermConstraints(t *testing.T) {
	order := &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.BTC, Term: instr.USD},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		Side:             side.BUY,
		Price:            100.0,
		Qty:              100.0,
		TimeInForce:      tif.GTC,
	}

	err := orders.ValidateLimitWithTermConstraints(order, 10.0, 20_000.0)
	if err != nil {
		t.Errorf("ValidateOrderWithTermConstraints() = %v; want nil", err)
	}

	err = orders.ValidateLimitWithTermConstraints(order, 10.0, 10_000.0)
	if err != nil {
		t.Errorf("ValidateOrderWithTermConstraints() = %v; want nil", err)
	}

	err = orders.ValidateLimitWithTermConstraints(order, 10.0, 9_999.0)
	if err == nil {
		t.Errorf("ValidateOrderWithTermConstraints() = nil; want error")
	}

	order.Qty = 0.001
	err = orders.ValidateLimitWithTermConstraints(order, 10.0, 1000.0)
	if err == nil {
		t.Errorf("ValidateOrderWithTermConstraints() = nil; want error")
	}

	order.Qty = 100.0
	order.Price = 0.0
	err = orders.ValidateLimitWithTermConstraints(order, 10.0, 1000.0)
	if err == nil {
		t.Errorf("ValidateOrderWithTermConstraints() = nil; want error")
	}

	order.Price = -1.0
	err = orders.ValidateLimitWithTermConstraints(order, 10.0, 1000.0)
	if err == nil {
		t.Errorf("ValidateOrderWithTermConstraints() = nil; want error")
	}
}

func TestValidateOrderWithRefData(t *testing.T) {
	order := &orders.Limit[instr.Spot]{}
	refData := &refdata.Composite{}

	// only check the error type
	if err := orders.ValidateLimitWithRefData(order, refData); err == nil {
		t.Errorf("Expected error")
	} else if err != refdata.NilRefDataError {
		t.Errorf("Wrong error type")
	}

	refData.PriceQty = struct{}{}
	refData.TradeSizeLimit = struct{}{}
	// only check the error type
	if err := orders.ValidateLimitWithRefData(order, refData); err == nil {
		t.Errorf("Expected error")
	} else if err != refdata.UnsupportedRefDataError {
		t.Errorf("Wrong error type: %v", err)
	}
}
