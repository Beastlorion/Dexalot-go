package inventory_test

import (
	"sync"
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/inventory"
)

func TestSimpleManager(t *testing.T) {
	inv := inventory.NewSimpleManager()
	inv.Set(inventory.Update{Asset: instr.BTC, Exchange: exchange.Binance, Qty: 100})

	if val, ok := inv.Get(instr.BTC, exchange.Binance); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 100 {
		t.Errorf("expected inventory of 100, got %f", val)
	}

	inv.Set(inventory.Update{Asset: instr.ETH, Exchange: exchange.Binance, Qty: 150})
	if val, ok := inv.Get(instr.ETH, exchange.Binance); !ok {
		t.Errorf("expected to find ETH, got %v", val)
	} else if val != 150 {
		t.Errorf("expected inventory of 150, got %f", val)
	}

	inv.Set(inventory.Update{Asset: instr.BTC, Exchange: exchange.Binance, Qty: 200})
	if val, ok := inv.Get(instr.BTC, exchange.Binance); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 200 {
		t.Errorf("expected inventory of 200, got %f", val)
	}

	inv.Update(inventory.Update{Asset: instr.BTC, Exchange: exchange.Binance, Qty: 100})
	if val, ok := inv.Get(instr.BTC, exchange.Binance); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 300 {
		t.Errorf("expected inventory of 200, got %f", val)
	}

	inv.Update(inventory.Update{Asset: instr.ETH, Exchange: exchange.Binance, Qty: -100})
	if val, ok := inv.Get(instr.ETH, exchange.Binance); !ok {
		t.Errorf("expected to find ETH, got %v", val)
	} else if val != 50 {
		t.Errorf("expected inventory of 50, got %f", val)
	}

	inv.Set(inventory.Update{Asset: instr.BTC, Exchange: exchange.Coinbase, Qty: 100})
	if val, ok := inv.Get(instr.BTC, exchange.Coinbase); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 100 {
		t.Errorf("expected inventory of 100, got %f", val)
	} else if otherVal, otherOk := inv.Get(instr.BTC, exchange.Binance); !otherOk {
		t.Errorf("expected to find BTC, got %v", otherVal)
	} else if otherVal != 300 {
		t.Errorf("expected inventory of 300, got %f", otherVal)
	}

	inv.Update(inventory.Update{Asset: instr.BTC, Exchange: exchange.Coinbase, Qty: -200})
	if val, ok := inv.Get(instr.BTC, exchange.Coinbase); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != -100 {
		t.Errorf("expected inventory of -100, got %f", val)
	} else if otherVal, otherOk := inv.Get(instr.BTC, exchange.Binance); !otherOk {
		t.Errorf("expected to find BTC, got %v", otherVal)
	} else if otherVal != 300 {
		t.Errorf("expected inventory of 300, got %f", otherVal)
	}

	inv.Update(inventory.Update{Asset: instr.ETH, Exchange: exchange.Coinbase, Qty: 500})
	if val, ok := inv.Get(instr.ETH, exchange.Coinbase); !ok {
		t.Errorf("expected to find ETH, got %v", val)
	} else if val != 500 {
		t.Errorf("expected inventory of 500, got %f", val)
	} else if otherVal, otherOk := inv.Get(instr.ETH, exchange.Binance); !otherOk {
		t.Errorf("expected to find ETH, got %v", otherVal)
	} else if otherVal != 50 {
		t.Errorf("expected inventory of 50, got %f", otherVal)
	}
}

func TestAtomicManager(t *testing.T) {
	inv := inventory.NewAtomicManager()
	inv.Set(inventory.Update{Asset: instr.BTC, Exchange: exchange.Binance, Qty: 100})

	if val, ok := inv.Get(instr.BTC, exchange.Binance); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 100 {
		t.Errorf("expected inventory of 100, got %f", val)
	}

	inv.Set(inventory.Update{Asset: instr.ETH, Exchange: exchange.Binance, Qty: 150})
	if val, ok := inv.Get(instr.ETH, exchange.Binance); !ok {
		t.Errorf("expected to find ETH, got %v", val)
	} else if val != 150 {
		t.Errorf("expected inventory of 150, got %f", val)
	}

	inv.Set(inventory.Update{Asset: instr.BTC, Exchange: exchange.Binance, Qty: 200})
	if val, ok := inv.Get(instr.BTC, exchange.Binance); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 200 {
		t.Errorf("expected inventory of 200, got %f", val)
	}

	inv.Update(inventory.Update{Asset: instr.BTC, Exchange: exchange.Binance, Qty: 100})
	if val, ok := inv.Get(instr.BTC, exchange.Binance); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 300 {
		t.Errorf("expected inventory of 200, got %f", val)
	}

	inv.Update(inventory.Update{Asset: instr.ETH, Exchange: exchange.Binance, Qty: -100})
	if val, ok := inv.Get(instr.ETH, exchange.Binance); !ok {
		t.Errorf("expected to find ETH, got %v", val)
	} else if val != 50 {
		t.Errorf("expected inventory of 50, got %f", val)
	}

	inv.Set(inventory.Update{Asset: instr.BTC, Exchange: exchange.Coinbase, Qty: 100})
	if val, ok := inv.Get(instr.BTC, exchange.Coinbase); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != 100 {
		t.Errorf("expected inventory of 100, got %f", val)
	} else if otherVal, otherOk := inv.Get(instr.BTC, exchange.Binance); !otherOk {
		t.Errorf("expected to find BTC, got %v", otherVal)
	} else if otherVal != 300 {
		t.Errorf("expected inventory of 300, got %f", otherVal)
	}

	inv.Update(inventory.Update{Asset: instr.BTC, Exchange: exchange.Coinbase, Qty: -200})
	if val, ok := inv.Get(instr.BTC, exchange.Coinbase); !ok {
		t.Errorf("expected to find BTC, got %v", val)
	} else if val != -100 {
		t.Errorf("expected inventory of -100, got %f", val)
	} else if otherVal, otherOk := inv.Get(instr.BTC, exchange.Binance); !otherOk {
		t.Errorf("expected to find BTC, got %v", otherVal)
	} else if otherVal != 300 {
		t.Errorf("expected inventory of 300, got %f", otherVal)
	}

	inv.Update(inventory.Update{Asset: instr.ETH, Exchange: exchange.Coinbase, Qty: 500})
	if val, ok := inv.Get(instr.ETH, exchange.Coinbase); !ok {
		t.Errorf("expected to find ETH, got %v", val)
	} else if val != 500 {
		t.Errorf("expected inventory of 500, got %f", val)
	} else if otherVal, otherOk := inv.Get(instr.ETH, exchange.Binance); !otherOk {
		t.Errorf("expected to find ETH, got %v", otherVal)
	} else if otherVal != 50 {
		t.Errorf("expected inventory of 50, got %f", otherVal)
	}

	var wg sync.WaitGroup
	n := 100
	iterations := 10000

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				inv.Update(inventory.Update{Asset: instr.ETH, Exchange: exchange.Dexalot, Qty: 1})
			}
		}()
	}

	wg.Wait()

	i, _ := inv.Get(instr.ETH, exchange.Dexalot)

	expected := float64(n * iterations)
	if i != expected {
		t.Errorf("expected inventory to be %v, but got %v", expected, i)
	}
}
