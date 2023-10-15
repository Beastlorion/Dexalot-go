package maker_test

import (
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordstatus"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordtypes"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/engine/maker"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/report"
	"github.com/stretchr/testify/assert"
)

func TestMakerLayerController_HasMetMinimumQuoteTime(t *testing.T) {
	spot := instr.Spot{Base: instr.BTC, Term: instr.USD}
	ex := exchange.Coinbase

	controller := maker.NewMakerLayerController(spot, ex)

	minQuoteTime := 100 * time.Millisecond
	now := time.Now()

	order := &orders.Limit[instr.Spot]{
		Instrument:       spot,
		ExchangeMetaData: &exchange.MetaData{Exchange: ex},
		Side:             side.BUY,
		Price:            10_000.0,
		Qty:              1.0,
		TimeInForce:      tif.GTC,
		TransactTime:     time.Now().Add(-time.Second),
	}

	assert.True(t, controller.HasMetMinimumQuoteTime(order, minQuoteTime, now))

	order.TransactTime = now

	assert.False(t, controller.HasMetMinimumQuoteTime(order, minQuoteTime, now))

	now = now.Add(100 * time.Millisecond)

	assert.True(t, controller.HasMetMinimumQuoteTime(order, minQuoteTime, now))

}

func TestMakerLayerController_GetLayerStatus(t *testing.T) {
	instrument := instr.Spot{Base: instr.BTC, Term: instr.USD}
	ex := exchange.Coinbase
	refData := refdata.Composite{
		ReplenishmentRate: 10 * time.Millisecond,
	}
	now := time.Now()
	controller := maker.NewMakerLayerController(instrument, ex)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "1",
		Side:          side.BUY,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "2",
		Side:          side.SELL,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Canceled,
		ClientOrderID: "1",
		Side:          side.BUY,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Canceled,
		ClientOrderID: "2",
		Side:          side.SELL,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.PO,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "1",
		Side:          side.BUY,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.PO,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "2",
		Side:          side.SELL,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.PO,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "1",
		Side:          side.BUY,
		TransactTime:  now,
	})

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.PO,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "2",
		Side:          side.SELL,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.PO,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "1",
		Side:          side.BUY,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "3",
		Side:          side.BUY,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "4",
		Side:          side.BUY,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "5",
		Side:          side.BUY,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "3",
		Side:          side.BUY,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "4",
		Side:          side.BUY,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "5",
		Side:          side.BUY,
		TransactTime:  now,
	})

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "6",
		Side:          side.SELL,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "7",
		Side:          side.SELL,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.New,
		ClientOrderID: "8",
		Side:          side.SELL,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "6",
		Side:          side.SELL,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "7",
		Side:          side.SELL,
		TransactTime:  now,
	})

	controller.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instrument,
		Exchange:      ex,
		OrderType:     ordtypes.Limit,
		TimeInForce:   tif.GTC,
		OrderStatus:   ordstatus.Filled,
		ClientOrderID: "8",
		Side:          side.SELL,
		TransactTime:  now,
	})

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.NotReplenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Replenishable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))

	now = now.Add(5 * time.Millisecond)

	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(0, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(1, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(2, side.SELL, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.BUY, &refData, now))
	assert.Equal(t, maker.Reconcilable, controller.GetLayerStatus(3, side.SELL, &refData, now))
}
