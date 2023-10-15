package refdata_test

import (
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
	"github.com/stretchr/testify/assert"
)

func TestRefDataManager(t *testing.T) {
	exchange1 := exchange.Coinbase
	spot1 := instr.Spot{Base: instr.BTC, Term: instr.USD}

	precisionRefData := refdata.Precision{
		PricePrecision: 2,
		QtyPrecision:   3,
	}

	baseLimitRefData := refdata.BaseTradeSizeLimit{
		BaseMinTradeSize: 0.01,
		BaseMaxTradeSize: 1_000.0,
	}

	refData1 := &refdata.Composite{
		PriceQty:          precisionRefData,
		TradeSizeLimit:    baseLimitRefData,
		MinQuoteTime:      1 * time.Second,
		ReplenishmentRate: 5 * time.Second,
		UsePostOnly:       false,
	}

	exchange2 := exchange.Binance
	spot2 := instr.Spot{Base: instr.ETH, Term: instr.USD}

	incrementRefData := refdata.Increment{
		PriceIncrement: 0.01,
		QtyIncrement:   0.001,
	}

	termLimitRefData := refdata.TermTradeSizeLimit{
		TermMinTradeSize: 5.0,
		TermMaxTradeSize: 100_000.0,
	}

	refData2 := &refdata.Composite{
		PriceQty:          incrementRefData,
		TradeSizeLimit:    termLimitRefData,
		MinQuoteTime:      1 * time.Second,
		ReplenishmentRate: 5 * time.Second,
		UsePostOnly:       false,
	}

	manager := refdata.NewManager[instr.Spot]()

	manager.Register(exchange1, spot1, refData1)
	manager.Register(exchange2, spot2, refData2)

	val, err := manager.Get(exchange1, spot1)
	assert.NoError(t, err)
	assert.Equal(t, refData1, val)

	val, err = manager.Get(exchange2, spot2)
	assert.NoError(t, err)
	assert.Equal(t, refData2, val)

	val, err = manager.Get(exchange1, spot2)
	assert.Error(t, err)
	assert.IsType(t, refdata.ReferenceDataNotFoundError, err)
	assert.Nil(t, val)

	val, err = manager.Get(exchange2, spot1)
	assert.Error(t, err)
	assert.IsType(t, refdata.ReferenceDataNotFoundError, err)
	assert.Nil(t, val)

	val, err = manager.Get(exchange.Dexalot, instr.Spot{Base: instr.AVAX, Term: instr.USD})
	assert.Error(t, err)
	assert.IsType(t, refdata.ReferenceDataNotFoundError, err)
	assert.Nil(t, val)
}
