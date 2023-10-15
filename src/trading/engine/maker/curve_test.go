package maker_test

import (
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/engine/maker"
	"github.com/stretchr/testify/assert"
)

const epsilon float64 = 0.0000001

func TestLiquidityCurveModel(t *testing.T) {
	instrument := instr.Spot{Base: instr.BTC, Term: instr.USD}
	ex := exchange.Binance
	modelFunc := func(totalQty float64) float64 {
		return totalQty
	}

	spreadModel := maker.NewStaticSymmetricCurveModel(modelFunc)

	lcm := maker.NewFunctionalLiqudityCurve(spreadModel, instrument, ex)
	now := time.Now()

	bids := lcm.GenerateBids(100.0, now)
	assert.Equal(t, 0, len(bids))

	offers := lcm.GenerateOffers(100.0, now)
	assert.Equal(t, 0, len(offers))

	lcm.SetLevelQtysInBase([]float64{100.0, 400.0, 250.0, 500.0})

	mid := 100.0

	bids = lcm.GenerateBids(mid, now)
	assert.Equal(t, 4, len(bids))

	assert.Equal(t, instrument, bids[0].Instrument)
	assert.Equal(t, ex, bids[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 99.0, bids[0].Price, epsilon)
	assert.Equal(t, 100.0, bids[0].Qty)
	assert.Equal(t, side.BUY, bids[0].Side)

	assert.Equal(t, instrument, bids[1].Instrument)
	assert.Equal(t, ex, bids[1].ExchangeMetaData.Exchange)
	assert.InDelta(t, 95.0, bids[1].Price, epsilon)
	assert.Equal(t, 400.0, bids[1].Qty)
	assert.Equal(t, side.BUY, bids[1].Side)

	assert.Equal(t, instrument, bids[2].Instrument)
	assert.Equal(t, ex, bids[2].ExchangeMetaData.Exchange)
	assert.InDelta(t, 92.5, bids[2].Price, epsilon)
	assert.Equal(t, 250.0, bids[2].Qty)
	assert.Equal(t, side.BUY, bids[2].Side)

	assert.Equal(t, instrument, bids[3].Instrument)
	assert.Equal(t, ex, bids[3].ExchangeMetaData.Exchange)
	assert.InDelta(t, 87.5, bids[3].Price, epsilon)
	assert.Equal(t, 500.0, bids[3].Qty)
	assert.Equal(t, side.BUY, bids[3].Side)

	offers = lcm.GenerateOffers(mid, now)
	assert.Equal(t, 4, len(offers))

	assert.Equal(t, instrument, offers[0].Instrument)
	assert.Equal(t, ex, offers[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 101.0, offers[0].Price, epsilon)
	assert.Equal(t, 100.0, offers[0].Qty)
	assert.Equal(t, side.SELL, offers[0].Side)

	assert.Equal(t, instrument, offers[1].Instrument)
	assert.Equal(t, ex, offers[1].ExchangeMetaData.Exchange)
	assert.InDelta(t, 105.0, offers[1].Price, epsilon)
	assert.Equal(t, 400.0, offers[1].Qty)
	assert.Equal(t, side.SELL, offers[1].Side)

	assert.Equal(t, instrument, offers[2].Instrument)
	assert.Equal(t, ex, offers[2].ExchangeMetaData.Exchange)
	assert.InDelta(t, 107.5, offers[2].Price, epsilon)
	assert.Equal(t, 250.0, offers[2].Qty)
	assert.Equal(t, side.SELL, offers[2].Side)

	assert.Equal(t, instrument, offers[3].Instrument)
	assert.Equal(t, ex, offers[3].ExchangeMetaData.Exchange)
	assert.InDelta(t, 112.5, offers[3].Price, epsilon)
	assert.Equal(t, 500.0, offers[3].Qty)
	assert.Equal(t, side.SELL, offers[3].Side)

	lcm.SetLevelQtysInBase([]float64{250.0, 500.0})

	bids = lcm.GenerateBids(mid, now)
	assert.Equal(t, 2, len(bids))

	assert.Equal(t, instrument, bids[0].Instrument)
	assert.Equal(t, ex, bids[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 97.5, bids[0].Price, epsilon)
	assert.Equal(t, 250.0, bids[0].Qty)
	assert.Equal(t, side.BUY, bids[0].Side)

	assert.Equal(t, instrument, bids[1].Instrument)
	assert.Equal(t, ex, bids[1].ExchangeMetaData.Exchange)
	assert.InDelta(t, 92.5, bids[1].Price, epsilon)
	assert.Equal(t, 500.0, bids[1].Qty)
	assert.Equal(t, side.BUY, bids[1].Side)

	offers = lcm.GenerateOffers(mid, now)
	assert.Equal(t, 2, len(offers))

	assert.Equal(t, instrument, offers[0].Instrument)
	assert.Equal(t, ex, offers[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 102.5, offers[0].Price, epsilon)
	assert.Equal(t, 250.0, offers[0].Qty)
	assert.Equal(t, side.SELL, offers[0].Side)

	assert.Equal(t, instrument, offers[1].Instrument)
	assert.Equal(t, ex, offers[1].ExchangeMetaData.Exchange)
	assert.InDelta(t, 107.5, offers[1].Price, epsilon)
	assert.Equal(t, 500.0, offers[1].Qty)
	assert.Equal(t, side.SELL, offers[1].Side)
}

func TestSingleLayerLiquidityCurveModel(t *testing.T) {
	instrument := instr.Spot{Base: instr.BTC, Term: instr.USD}
	ex := exchange.Coinbase
	model := maker.NewSingleLayerSpreadModel()
	liquidityCurve := maker.NewSingleLayerLiquidityCurve(model, instrument, ex)

	model.SetSpreadBps(10.0)
	liquidityCurve.SetBaseInventory(1.0)
	liquidityCurve.SetTermInventory(10_000.0)

	mid := 10_000.0
	now := time.Now()

	bids := liquidityCurve.GenerateBids(mid, now)
	assert.Equal(t, 1, len(bids))

	assert.Equal(t, instrument, bids[0].Instrument)
	assert.Equal(t, ex, bids[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 9_990.0, bids[0].Price, epsilon)
	assert.InDelta(t, 1.001, bids[0].Qty, 1e-5)
	assert.Equal(t, side.BUY, bids[0].Side)

	offers := liquidityCurve.GenerateOffers(mid, now)
	assert.Equal(t, 1, len(offers))

	assert.Equal(t, instrument, offers[0].Instrument)
	assert.Equal(t, ex, offers[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 10_010.0, offers[0].Price, epsilon)
	assert.InDelta(t, 1.0, offers[0].Qty, epsilon)
	assert.Equal(t, side.SELL, offers[0].Side)

	mid = 11_000.0

	bids = liquidityCurve.GenerateBids(mid, now)
	assert.Equal(t, 1, len(bids))

	assert.Equal(t, instrument, bids[0].Instrument)
	assert.Equal(t, ex, bids[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 10_989.0, bids[0].Price, epsilon)
	assert.InDelta(t, 0.91, bids[0].Qty, 1e-5)
	assert.Equal(t, side.BUY, bids[0].Side)

	offers = liquidityCurve.GenerateOffers(mid, now)
	assert.Equal(t, 1, len(offers))

	assert.Equal(t, instrument, offers[0].Instrument)
	assert.Equal(t, ex, offers[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 11_011.0, offers[0].Price, epsilon)
	assert.InDelta(t, 1.0, offers[0].Qty, epsilon)
	assert.Equal(t, side.SELL, offers[0].Side)

	mid = 10_000.0
	liquidityCurve.SetBaseInventory(0.0)
	liquidityCurve.SetTermInventory(0.0)

	bids = liquidityCurve.GenerateBids(mid, now)
	assert.Equal(t, 1, len(bids))

	assert.Equal(t, instrument, bids[0].Instrument)
	assert.Equal(t, ex, bids[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 9_990.0, bids[0].Price, epsilon)
	assert.InDelta(t, 0.0, bids[0].Qty, epsilon)
	assert.Equal(t, side.BUY, bids[0].Side)

	offers = liquidityCurve.GenerateOffers(mid, now)
	assert.Equal(t, 1, len(offers))

	assert.Equal(t, instrument, offers[0].Instrument)
	assert.Equal(t, ex, offers[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 10_010.0, offers[0].Price, epsilon)
	assert.InDelta(t, 0.0, offers[0].Qty, epsilon)
	assert.Equal(t, side.SELL, offers[0].Side)

	liquidityCurve.SetBaseInventory(1.0)
	liquidityCurve.SetTermInventory(10_000.0)
	model.SetSpreadBps(0.0)

	bids = liquidityCurve.GenerateBids(mid, now)
	assert.Equal(t, 1, len(bids))

	assert.Equal(t, instrument, bids[0].Instrument)
	assert.Equal(t, ex, bids[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 10_000.0, bids[0].Price, epsilon)
	assert.InDelta(t, 1.0, bids[0].Qty, 1e-5)
	assert.Equal(t, side.BUY, bids[0].Side)

	offers = liquidityCurve.GenerateOffers(mid, now)
	assert.Equal(t, 1, len(offers))

	assert.Equal(t, instrument, offers[0].Instrument)
	assert.Equal(t, ex, offers[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 10_000.0, offers[0].Price, epsilon)
	assert.InDelta(t, 1.0, offers[0].Qty, epsilon)
	assert.Equal(t, side.SELL, offers[0].Side)

	model.SetSpreadBps(50.0)

	bids = liquidityCurve.GenerateBids(mid, now)
	assert.Equal(t, 1, len(bids))

	assert.Equal(t, instrument, bids[0].Instrument)
	assert.Equal(t, ex, bids[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 9_950.0, bids[0].Price, epsilon)
	assert.InDelta(t, 1.005025, bids[0].Qty, 1e-5)
	assert.Equal(t, side.BUY, bids[0].Side)

	offers = liquidityCurve.GenerateOffers(mid, now)
	assert.Equal(t, 1, len(offers))

	assert.Equal(t, instrument, offers[0].Instrument)
	assert.Equal(t, ex, offers[0].ExchangeMetaData.Exchange)
	assert.InDelta(t, 10_050.0, offers[0].Price, epsilon)
	assert.InDelta(t, 1.0, offers[0].Qty, epsilon)
	assert.Equal(t, side.SELL, offers[0].Side)
}
