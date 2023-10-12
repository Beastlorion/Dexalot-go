package price_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/price"
	"github.com/stretchr/testify/assert"
)

const epsilon = 1e-12

func TestSkewPrice(t *testing.T) {
	px := 100.0
	skewBps := 100.0

	skewedPrice := price.SkewPrice(px, skewBps)
	assert.InDelta(t, 101.0, skewedPrice, epsilon)

	skewedPrice = price.SkewPrice(px, -skewBps)

	assert.InDelta(t, 99.0, skewedPrice, epsilon)

	skewBps = 0.0

	skewedPrice = price.SkewPrice(px, skewBps)
	assert.InDelta(t, 100.0, skewedPrice, epsilon)

	skewedPrice = price.SkewPrice(px, -skewBps)
	assert.InDelta(t, 100.0, skewedPrice, epsilon)

	skewBps = 5.0

	skewedPrice = price.SkewPrice(px, skewBps)
	assert.InDelta(t, 100.05, skewedPrice, epsilon)

	skewedPrice = price.SkewPrice(px, -skewBps)
	assert.InDelta(t, 99.95, skewedPrice, epsilon)
}

func TestAddSpread(t *testing.T) {
	px := 100.0
	spread := 100.0

	upPrice := price.AddSpread(px, spread, price.Up)
	assert.InDelta(t, 101.0, upPrice, epsilon)

	downPrice := price.AddSpread(px, spread, price.Down)
	assert.InDelta(t, 99.0, downPrice, epsilon)

	spread = 0.0

	upPrice = price.AddSpread(px, spread, price.Up)
	assert.InDelta(t, 100.0, upPrice, epsilon)

	downPrice = price.AddSpread(px, spread, price.Down)
	assert.InDelta(t, 100.0, downPrice, epsilon)

	spread = 5.0

	upPrice = price.AddSpread(px, spread, price.Up)
	assert.InDelta(t, 100.05, upPrice, epsilon)

	downPrice = price.AddSpread(px, spread, price.Down)
	assert.InDelta(t, 99.95, downPrice, epsilon)

	spread = -5.0

	upPrice = price.AddSpread(px, spread, price.Up)
	assert.InDelta(t, 99.95, upPrice, epsilon)

	downPrice = price.AddSpread(px, spread, price.Down)
	assert.InDelta(t, 100.05, downPrice, epsilon)

	assert.Panics(t, func() { price.AddSpread(px, spread, 0) })
}

func TestGetSpreadBps(t *testing.T) {
	mid := 100.0
	price2 := 100.0

	spreadBps := price.GetSpreadBps(mid, price2)
	assert.InDelta(t, 0.0, spreadBps, epsilon)

	price2 = 101.0

	spreadBps = price.GetSpreadBps(mid, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price2 = 99.0

	spreadBps = price.GetSpreadBps(mid, price2)
	assert.InDelta(t, -100.0, spreadBps, epsilon)

	price2 = 100.05

	spreadBps = price.GetSpreadBps(mid, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)

	price2 = 99.95

	spreadBps = price.GetSpreadBps(mid, price2)
	assert.InDelta(t, -5.0, spreadBps, epsilon)

}

func TestGetDistanceBps(t *testing.T) {
	mid := 100.0
	price1 := 100.0
	price2 := 100.0

	spreadBps := price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, 0.0, spreadBps, epsilon)

	price2 = 101.0

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price2 = 99.0

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, -100.0, spreadBps, epsilon)

	price2 = 100.05

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)

	price2 = 99.95

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, -5.0, spreadBps, epsilon)

	price1 = 101.0
	price2 = 100.0

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, -100.0, spreadBps, epsilon)

	price1 = 99.0

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price1 = 100.05

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, -5.0, spreadBps, epsilon)

	price1 = 99.95

	spreadBps = price.GetDistanceBps(mid, price1, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)

}

func TestGetAbsoluteSpreadBps(t *testing.T) {
	mid := 100.0
	price2 := 100.0

	spreadBps := price.GetAbsoluteSpreadBps(mid, price2)
	assert.InDelta(t, 0.0, spreadBps, epsilon)

	price2 = 101.0

	spreadBps = price.GetAbsoluteSpreadBps(mid, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price2 = 99.0

	spreadBps = price.GetAbsoluteSpreadBps(mid, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price2 = 100.05

	spreadBps = price.GetAbsoluteSpreadBps(mid, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)

	price2 = 99.95

	spreadBps = price.GetAbsoluteSpreadBps(mid, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)
}

func TestGetAbsoulteDistanceBps(t *testing.T) {
	mid := 100.0
	price1 := 100.0
	price2 := 100.0

	spreadBps := price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 0.0, spreadBps, epsilon)

	price2 = 101.0

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price2 = 99.0

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price2 = 100.05

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)

	price2 = 99.95

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)

	price1 = 101.0
	price2 = 100.0

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price1 = 99.0

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 100.0, spreadBps, epsilon)

	price1 = 100.05

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)

	price1 = 99.95

	spreadBps = price.GetAbsoluteDistanceBps(mid, price1, price2)
	assert.InDelta(t, 5.0, spreadBps, epsilon)
}
