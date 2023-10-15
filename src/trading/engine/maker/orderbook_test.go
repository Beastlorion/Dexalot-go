package maker_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordstatus"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordtypes"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/engine/maker"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/report"
	"github.com/stretchr/testify/assert"
)

func TestMakerOrderBook(t *testing.T) {
	makerOrderBook := maker.NewOrderBook(instr.Spot{Base: instr.AVAX, Term: instr.USDC}, exchange.Coinbase)

	assert.Zero(t, makerOrderBook.Bids.Size())
	assert.Zero(t, makerOrderBook.Offers.Size())

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument: instr.Spot{Base: instr.BTC, Term: instr.USDC},
		Exchange:   exchange.Coinbase,
		Price:      10_000.0,
		Side:       side.BUY,
	})

	assert.Zero(t, makerOrderBook.Bids.Size())
	assert.Zero(t, makerOrderBook.Offers.Size())

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument: instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:   exchange.Dexalot,
		Price:      10_000.0,
		Side:       side.BUY,
	})

	assert.Zero(t, makerOrderBook.Bids.Size())
	assert.Zero(t, makerOrderBook.Offers.Size())

	order := &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "1",
		Price:            10_000.0,
		Qty:              1.0,
		Side:             side.BUY,
		TimeInForce:      tif.GTC,
	}

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         10_000.0,
		Qty:           1.0,
		Side:          side.BUY,
		ClientOrderID: "1",
		OrderID:       "1",
		OrderType:     ordtypes.Limit,
		OrderStatus:   ordstatus.New,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 1, makerOrderBook.Bids.Size())
	assert.Zero(t, makerOrderBook.Offers.Size())
	if order, ok := makerOrderBook.Bids.Get(maker.Layer{Price: order.Price, ClientOrderID: "1"}); !ok {
		t.Fatalf("Expected to find order at 10_000.0")
	} else if order.ClientOrderID != order.ClientOrderID {
		t.Errorf("Expected order to have ClientOrderID 1, got %s", order.ClientOrderID)
	} else if *order.ExchangeMetaData.OrderID != "1" {
		t.Errorf("Expected order to have OrderID 1, got %s", *order.ExchangeMetaData.OrderID)
	} else if order.Qty != 1.0 {
		t.Errorf("Expected order to have Qty 1.0, got %f", order.Qty)
	} else if order.TimeInForce != tif.GTC {
		t.Errorf("Expected order to have TimeInForce GTC, got %d", order.TimeInForce)
	} else if order.ExchangeMetaData.Exchange != exchange.Coinbase {
		t.Errorf("Expected order to have ExchangeMetaData.Exchange.Name COINBASE, got %s", order.ExchangeMetaData.Exchange)
	} else if order.Instrument.String() != "AVAX-USDC" {
		t.Errorf("Expected order to have Instrument AVAX-USDC, got %s", order.Instrument.String())
	}

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         11_000.0,
		Side:          side.BUY,
		ClientOrderID: "2",
		OrderID:       "2",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		Qty:           1.4,
		TimeInForce:   tif.IOC,
	})

	assert.Equal(t, 1, makerOrderBook.Bids.Size())
	assert.Zero(t, makerOrderBook.Offers.Size())
	if order, ok := makerOrderBook.Bids.Get(maker.Layer{Price: order.Price, ClientOrderID: "2"}); ok {
		t.Fatalf("Expected not to find order at 11_000.0, got %+v", order)
	}

	order = &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "3",
		Price:            12_000.0,
		Qty:              2.0,
		Side:             side.SELL,
		TimeInForce:      tif.PO,
	}
	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Qty:           2.0,
		Side:          side.SELL,
		ClientOrderID: "3",
		OrderID:       "3",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		TimeInForce:   order.TimeInForce,
	})

	assert.Equal(t, 1, makerOrderBook.Bids.Size())
	assert.Equal(t, 1, makerOrderBook.Offers.Size())
	if order, ok := makerOrderBook.Offers.Get(maker.Layer{Price: order.Price, ClientOrderID: "3"}); !ok {
		t.Fatalf("Expected to find order at 12_000.0")
	} else if order.ClientOrderID != order.ClientOrderID {
		t.Errorf("Expected order to have ClientOrderID 3, got %s", order.ClientOrderID)
	} else if *order.ExchangeMetaData.OrderID != "3" {
		t.Errorf("Expected order to have OrderID 3, got %s", *order.ExchangeMetaData.OrderID)
	} else if order.Qty != 2.0 {
		t.Errorf("Expected order to have Qty 2.0, got %f", order.Qty)
	} else if order.TimeInForce != tif.PO {
		t.Errorf("Expected order to have TimeInForce PO, got %d", order.TimeInForce)
	} else if order.ExchangeMetaData.Exchange != exchange.Coinbase {
		t.Errorf("Expected order to have ExchangeMetaData.Exchange.Name COINBASE, got %s", order.ExchangeMetaData.Exchange)
	} else if order.Instrument.String() != "AVAX-USDC" {
		t.Errorf("Expected order to have Instrument AVAX-USDC, got %s", order.Instrument.String())
	}

	order = &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "4",
		Price:            12_500.0,
		Qty:              1.0,
		Side:             side.SELL,
		TimeInForce:      tif.GTC,
	}
	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_500.0,
		Qty:           1.0,
		Side:          side.SELL,
		ClientOrderID: "4",
		OrderID:       "4",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		TimeInForce:   order.TimeInForce,
	})

	order = &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "5",
		Price:            11_500.0,
		Qty:              1.0,
		Side:             side.SELL,
		TimeInForce:      tif.GTC,
	}
	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         11_500.0,
		Side:          side.SELL,
		ClientOrderID: "5",
		OrderID:       "5",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   order.TimeInForce,
	})

	offers := makerOrderBook.Offers.Values()
	assert.Equal(t, 3, len(offers))
	assert.Equal(t, 11_500.0, offers[0].Price)
	assert.Equal(t, 12_000.0, offers[1].Price)
	assert.Equal(t, 12_500.0, offers[2].Price)
	assert.Equal(t, "5", *offers[0].ExchangeMetaData.OrderID)
	assert.Equal(t, "3", *offers[1].ExchangeMetaData.OrderID)
	assert.Equal(t, "4", *offers[2].ExchangeMetaData.OrderID)
	assert.Equal(t, "5", offers[0].ClientOrderID)
	assert.Equal(t, "3", offers[1].ClientOrderID)
	assert.Equal(t, "4", offers[2].ClientOrderID)

	iter := makerOrderBook.Offers.Iterator()
	iter.Next()
	assert.Equal(t, maker.Layer{Price: 11_500.0, ClientOrderID: "5"}, iter.Key())
	assert.Equal(t, "5", *iter.Value().ExchangeMetaData.OrderID)
	iter.Next()
	assert.Equal(t, maker.Layer{Price: 12_000.0, ClientOrderID: "3"}, iter.Key())
	assert.Equal(t, "3", *iter.Value().ExchangeMetaData.OrderID)
	iter.Next()
	assert.Equal(t, maker.Layer{Price: 12_500.0, ClientOrderID: "4"}, iter.Key())
	assert.Equal(t, "4", *iter.Value().ExchangeMetaData.OrderID)

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         11_500.0,
		Side:          side.SELL,
		ClientOrderID: "5",
		OrderID:       "5",
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 2, makerOrderBook.Offers.Size())
	_, ok := makerOrderBook.Offers.Get(maker.Layer{Price: 12_500.0, ClientOrderID: "4"})
	if !ok {
		t.Fatalf("Expected to find order at 12_500.0")
	}
	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_500.0,
		Side:          side.SELL,
		ClientOrderID: "4",
		OrderID:       "4",
		OrderStatus:   ordstatus.Canceled,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 1, makerOrderBook.Offers.Size())
	if order, ok := makerOrderBook.Offers.Get(maker.Layer{Price: 12_500.0, ClientOrderID: "4"}); ok {
		t.Fatalf("Expected not to find order at 12_500.0, got %+v", order)
	}

	_, ok = makerOrderBook.Offers.Get(maker.Layer{Price: 12_000.0, ClientOrderID: "3"})
	if !ok {
		t.Fatalf("Expected to find order at 12_000.0")
	}
	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.SELL,
		ClientOrderID: "3",
		OrderID:       "3",
		OrderStatus:   ordstatus.DoneForDay,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 0, makerOrderBook.Offers.Size())
	if order, ok := makerOrderBook.Offers.Get(maker.Layer{Price: 12_000.0, ClientOrderID: "3"}); ok {
		t.Fatalf("Expected not to find order at 12_000.0, got %+v", order)
	}

	assert.Equal(t, 1, makerOrderBook.Bids.Size())

	_, ok = makerOrderBook.Bids.Get(maker.Layer{Price: 10_000.0, ClientOrderID: "1"})
	if !ok {
		t.Fatalf("Expected to find order at 10_000.0")
	}
	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         10_000.0,
		Side:          side.BUY,
		ClientOrderID: "1",
		OrderID:       "1",
		OrderStatus:   ordstatus.Expired,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 0, makerOrderBook.Bids.Size())
	if order, ok := makerOrderBook.Bids.Get(maker.Layer{Price: 10_000.0, ClientOrderID: "1"}); ok {
		t.Fatalf("Expected not to find order at 10_000.0, got %+v", order)
	}

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.BUY,
		ClientOrderID: "10",
		OrderID:       "10",
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 0, makerOrderBook.Bids.Size())
	if order, ok := makerOrderBook.Bids.Get(maker.Layer{Price: 12_000.0, ClientOrderID: "10"}); ok {
		t.Fatalf("Expected not to find order at 12_000.0, got %+v", order)
	}

	collidingBid1 := &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "11",
		Price:            12_000.0,
		Qty:              1.0,
		Side:             side.BUY,
		TimeInForce:      tif.GTC,
	}

	collidingBid2 := &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "12",
		Price:            12_000.0,
		Qty:              3.0,
		Side:             side.BUY,
		TimeInForce:      tif.GTC,
	}

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.BUY,
		ClientOrderID: "11",
		OrderID:       "11",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.BUY,
		ClientOrderID: "12",
		OrderID:       "12",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		Qty:           3.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 2, makerOrderBook.Bids.Size())

	iter = makerOrderBook.Bids.Iterator()
	assert.True(t, iter.Next())
	assert.Equal(t, collidingBid1.ClientOrderID, iter.Value().ClientOrderID)
	assert.True(t, iter.Next())
	assert.Equal(t, collidingBid2.ClientOrderID, iter.Value().ClientOrderID)

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.BUY,
		ClientOrderID: "11",
		OrderID:       "11",
		OrderStatus:   ordstatus.Canceled,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 1, makerOrderBook.Bids.Size())

	iter = makerOrderBook.Bids.Iterator()
	assert.True(t, iter.Next())
	assert.Equal(t, collidingBid2.ClientOrderID, iter.Value().ClientOrderID)

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.BUY,
		ClientOrderID: "12",
		OrderID:       "12",
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Qty:           3.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 0, makerOrderBook.Bids.Size())

	collidingOffer1 := &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "13",
		Price:            12_000.0,
		Qty:              1.0,
		Side:             side.SELL,
		TimeInForce:      tif.GTC,
	}

	collidingOffer2 := &orders.Limit[instr.Spot]{
		Instrument:       instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		ExchangeMetaData: &exchange.MetaData{Exchange: exchange.Coinbase},
		ClientOrderID:    "14",
		Price:            12_000.0,
		Qty:              3.0,
		Side:             side.SELL,
		TimeInForce:      tif.GTC,
	}

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.SELL,
		ClientOrderID: "13",
		OrderID:       "13",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.SELL,
		ClientOrderID: "14",
		OrderID:       "14",
		OrderStatus:   ordstatus.New,
		OrderType:     ordtypes.Limit,
		Qty:           3.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 2, makerOrderBook.Offers.Size())

	iter = makerOrderBook.Offers.Iterator()
	assert.True(t, iter.Next())
	assert.Equal(t, collidingOffer1.ClientOrderID, iter.Value().ClientOrderID)
	assert.True(t, iter.Next())
	assert.Equal(t, collidingOffer2.ClientOrderID, iter.Value().ClientOrderID)

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.SELL,
		ClientOrderID: "14",
		OrderID:       "14",
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Qty:           3.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 1, makerOrderBook.Offers.Size())

	iter = makerOrderBook.Offers.Iterator()
	assert.True(t, iter.Next())
	assert.Equal(t, collidingOffer1.ClientOrderID, iter.Value().ClientOrderID)
	assert.False(t, iter.Next())

	makerOrderBook.ApplyExecutionReport(&report.ExecutionReport[instr.Spot]{
		Instrument:    instr.Spot{Base: instr.AVAX, Term: instr.USDC},
		Exchange:      exchange.Coinbase,
		Price:         12_000.0,
		Side:          side.SELL,
		ClientOrderID: "13",
		OrderID:       "13",
		OrderStatus:   ordstatus.Canceled,
		OrderType:     ordtypes.Limit,
		Qty:           1.0,
		TimeInForce:   tif.GTC,
	})

	assert.Equal(t, 0, makerOrderBook.Offers.Size())
}
