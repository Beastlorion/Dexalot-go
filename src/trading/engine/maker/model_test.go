package maker_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/engine/maker"
	"github.com/stretchr/testify/assert"
)

func TestSingleLayerSpreadModel(t *testing.T) {
	model := maker.NewSingleLayerSpreadModel()

	model.SetSpreadBps(50.0)

	assert.Equal(t, 50.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 50.0, model.OfferSpreadBps(100.0))
	assert.Equal(t, 50.0, model.BidSpreadBps(1000.0))
	assert.Equal(t, 50.0, model.OfferSpreadBps(1000.0))
	assert.Equal(t, 50.0, model.BidSpreadBps(0.0))
	assert.Equal(t, 50.0, model.OfferSpreadBps(0.0))

	model.SetSpreadBps(100.0)

	assert.Equal(t, 100.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 100.0, model.OfferSpreadBps(100.0))
	assert.Equal(t, 100.0, model.BidSpreadBps(1000.0))
	assert.Equal(t, 100.0, model.OfferSpreadBps(1000.0))
	assert.Equal(t, 100.0, model.BidSpreadBps(0.0))
	assert.Equal(t, 100.0, model.OfferSpreadBps(0.0))

	model.SetSpreadBps(0.0)

	assert.Equal(t, 0.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 0.0, model.OfferSpreadBps(100.0))
	assert.Equal(t, 0.0, model.BidSpreadBps(1000.0))
	assert.Equal(t, 0.0, model.OfferSpreadBps(1000.0))
	assert.Equal(t, 0.0, model.BidSpreadBps(0.0))
	assert.Equal(t, 0.0, model.OfferSpreadBps(0.0))

	model.SetSpreadBps(-50.0)

	assert.Equal(t, -50.0, model.BidSpreadBps(100.0))
	assert.Equal(t, -50.0, model.OfferSpreadBps(100.0))
	assert.Equal(t, -50.0, model.BidSpreadBps(1000.0))
	assert.Equal(t, -50.0, model.OfferSpreadBps(1000.0))
	assert.Equal(t, -50.0, model.BidSpreadBps(0.0))
	assert.Equal(t, -50.0, model.OfferSpreadBps(0.0))

}

func TestWithBiasSpreadModel(t *testing.T) {
	model := maker.NewWithBiasSpreadModel(maker.NewSingleLayerSpreadModel())

	model.SetBidBiasBps(50.0)
	model.SetOfferBiasBps(25.0)

	assert.Equal(t, 50.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 25.0, model.OfferSpreadBps(100.0))

	model.SetBidBiasBps(0.0)
	model.SetOfferBiasBps(0.0)
	model.UnderlyingModel.SetSpreadBps(50.0)

	assert.Equal(t, 50.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 50.0, model.OfferSpreadBps(100.0))

	model.SetBidBiasBps(-50.0)
	model.SetOfferBiasBps(-25.0)

	assert.Equal(t, 0.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 25.0, model.OfferSpreadBps(100.0))

	model.SetBidBiasBps(40.0)
	model.SetOfferBiasBps(20.0)

	assert.Equal(t, 90.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 70.0, model.OfferSpreadBps(100.0))

}

func TestStaticSymmetricCurveModel(t *testing.T) {
	model := maker.NewStaticSymmetricCurveModel(
		func(totalQty float64) float64 {
			return totalQty // as bps
		},
	)

	assert.Equal(t, 100.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 1000.0, model.BidSpreadBps(1000.0))
	assert.Equal(t, 0.0, model.BidSpreadBps(0.0))
	assert.Equal(t, -100.0, model.BidSpreadBps(-100.0))

	assert.Equal(t, 100.0, model.OfferSpreadBps(100.0))
	assert.Equal(t, 1000.0, model.OfferSpreadBps(1000.0))
	assert.Equal(t, 0.0, model.OfferSpreadBps(0.0))
	assert.Equal(t, -100.0, model.OfferSpreadBps(-100.0))
}

func TestTwoFactorSymmetricCurveModel(t *testing.T) {
	model := maker.NewTwoFactorSymmetricCurveModel()

	model.SetInsideQty(10.0)
	model.SetInsideSpreadBps(5.0)
	model.SetQtySpreadMultiplierBps(1.0)

	assert.Equal(t, 5.0, model.BidSpreadBps(10.0))
	assert.Equal(t, 5.0, model.OfferSpreadBps(10.0))

	assert.Equal(t, 15.0, model.BidSpreadBps(20.0))
	assert.Equal(t, 15.0, model.OfferSpreadBps(20.0))

	assert.Equal(t, 25.0, model.BidSpreadBps(30.0))
	assert.Equal(t, 25.0, model.OfferSpreadBps(30.0))

	assert.Equal(t, 5.0, model.BidSpreadBps(5.0))
	assert.Equal(t, 5.0, model.OfferSpreadBps(5.0))

	model.SetInsideQty(0.0)

	assert.Equal(t, 15.0, model.BidSpreadBps(10.0))
	assert.Equal(t, 15.0, model.OfferSpreadBps(10.0))

	assert.Equal(t, 25.0, model.BidSpreadBps(20.0))
	assert.Equal(t, 25.0, model.OfferSpreadBps(20.0))

	model.SetInsideSpreadBps(50.0)

	assert.Equal(t, 60.0, model.BidSpreadBps(10.0))
	assert.Equal(t, 60.0, model.OfferSpreadBps(10.0))

	model.SetInsideQty(100.0)

	assert.Equal(t, 50.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 50.0, model.OfferSpreadBps(100.0))

	assert.Equal(t, 100.0, model.BidSpreadBps(150.0))
	assert.Equal(t, 100.0, model.OfferSpreadBps(150.0))

	model.SetQtySpreadMultiplierBps(2.0)

	assert.Equal(t, 50.0, model.BidSpreadBps(100.0))
	assert.Equal(t, 50.0, model.OfferSpreadBps(100.0))

	assert.Equal(t, 150.0, model.BidSpreadBps(150.0))
	assert.Equal(t, 150.0, model.OfferSpreadBps(150.0))
}
