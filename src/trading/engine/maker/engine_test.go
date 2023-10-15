package maker_test

import (
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/source"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/ticker"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/price"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordstatus"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordtypes"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/engine/maker"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/report"
	"github.com/stretchr/testify/assert"
)

type MockDispatcher[T instr.Instrument] struct {
	newOrders    []*orders.Limit[T]
	cancelOrders []*orders.Limit[T]
}

func (d *MockDispatcher[T]) DispatchLimitOrder(order *orders.Limit[T]) {
	d.newOrders = append(d.newOrders, order)
}

func (d *MockDispatcher[T]) DispatchCancelOrder(order *orders.Limit[T]) {
	d.cancelOrders = append(d.cancelOrders, order)
}

func (d *MockDispatcher[T]) Reset() {
	d.newOrders = nil
	d.cancelOrders = nil
}

func NewMockDispatcher[T instr.Instrument]() *MockDispatcher[T] {
	return &MockDispatcher[T]{
		newOrders:    make([]*orders.Limit[T], 0),
		cancelOrders: make([]*orders.Limit[T], 0),
	}
}

func TestMakerEngine_MovingPrices(t *testing.T) {
	spot := instr.Spot{Base: instr.BTC, Term: instr.USD}
	ex := exchange.Coinbase
	pqRefData := &refdata.Precision{
		PricePrecision: 2,
		QtyPrecision:   3,
	}
	tradeSizeLimit := &refdata.BaseTradeSizeLimit{
		BaseMinTradeSize:      0.001,
		BaseMaxTradeSize:      1000,
	}
	refData := refdata.Composite{
		PriceQty: pqRefData,
		TradeSizeLimit: tradeSizeLimit,
		MinQuoteTime:          1 * time.Millisecond,
		ReplenishmentRate:     10 * time.Millisecond,
		UsePostOnly:          false,
	}

	mid := 10_000.0
	baseInventory := 20.0
	termInventory := 200_000.0
	baseQtyLevels := []float64{1, 5, 10}
	now := time.Now()

	orderBook := maker.NewOrderBook(spot, ex)
	dispatcher := NewMockDispatcher[instr.Spot]()
	controller := maker.NewMakerLayerController(spot, ex)
	model := maker.NewFunctionalLiqudityCurve(
		maker.NewStaticSymmetricCurveModel(
			func(totalBaseQty float64) float64 {
				return totalBaseQty * 10.0 // 10 bps per 1 BTC
			}),
		spot,
		ex,
	)

	makerEngine := maker.NewEngine[instr.Spot](dispatcher, orderBook, controller, model)

	// Run Offers no levels set
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)
	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	// Run Bids no levels set
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)
	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	now = now.Add(5 * time.Millisecond)

	// Set levels
	model.SetLevelQtysInBase(baseQtyLevels)

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 3, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_010.0, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_060.0, dispatcher.newOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 10_160.0, dispatcher.newOrders[2].Price)
	assert.Equal(t, 10.0, dispatcher.newOrders[2].Qty)

	// Force confirmation of orders in the order book
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 3, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	assert.Equal(t, 9_990.0, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 9940.0, dispatcher.newOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 9840.0, dispatcher.newOrders[2].Price)
	assert.Equal(t, 10.0, dispatcher.newOrders[2].Qty)
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Run the same Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))

	// Run the same Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))

	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move the mid price up by 1 cent
	mid = 10_000.01

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move the mid down by 2 cents
	mid = 9_999.99

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move mid enough to displace the first level
	mid = 10_004.0

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 1, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_014.01, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_010.0, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()

	// Run bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 1, len(dispatcher.cancelOrders))
	assert.Equal(t, 9_993.99, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 9_990.0, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move mid down enough to displace the first level and more, but not the second level
	mid = 9_996.0

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 1, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_006.0, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_014.01, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()

	// Run bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 1, len(dispatcher.cancelOrders))
	assert.Equal(t, 9_986.0, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 9_993.99, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move mid back up to replace the first two levels
	mid = 10_025.0

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 2, len(dispatcher.newOrders))
	assert.Equal(t, 2, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_035.03, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_085.15, dispatcher.newOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 10_006.0, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	assert.Equal(t, 10_060.0, dispatcher.cancelOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.cancelOrders[1].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()

	// Run bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 2, len(dispatcher.newOrders))
	assert.Equal(t, 2, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_014.97, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 9_964.85, dispatcher.newOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 9_986.0, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	assert.Equal(t, 9_940.0, dispatcher.cancelOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.cancelOrders[1].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move bid up further to replace the third level
	mid = 10_080.0

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 3, len(dispatcher.newOrders))
	assert.Equal(t, 3, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_090.08, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_140.48, dispatcher.newOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 10_241.29, dispatcher.newOrders[2].Price)
	assert.Equal(t, 10.0, dispatcher.newOrders[2].Qty)
	assert.Equal(t, 10_035.03, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	assert.Equal(t, 10_085.15, dispatcher.cancelOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.cancelOrders[1].Qty)
	assert.Equal(t, 10_160.0, dispatcher.cancelOrders[2].Price)
	assert.Equal(t, 10.0, dispatcher.cancelOrders[2].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 3, len(dispatcher.newOrders))
	assert.Equal(t, 3, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_069.92, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_019.52, dispatcher.newOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 9_918.71, dispatcher.newOrders[2].Price)
	assert.Equal(t, 10.0, dispatcher.newOrders[2].Qty)
	assert.Equal(t, 10_014.97, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	assert.Equal(t, 9_964.85, dispatcher.cancelOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.cancelOrders[1].Qty)
	assert.Equal(t, 9_840.0, dispatcher.cancelOrders[2].Price)
	assert.Equal(t, 10.0, dispatcher.cancelOrders[2].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move the bid so high that we are breaching the term inventory on the bid
	mid = 13_100.0

	// Skip the Offers

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 3, len(dispatcher.newOrders))
	assert.Equal(t, 3, len(dispatcher.cancelOrders))
	assert.Equal(t, 13_086.90, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 13_021.40, dispatcher.newOrders[1].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 12_890.40, dispatcher.newOrders[2].Price)
	assert.Equal(t, 9.449, dispatcher.newOrders[2].Qty)

	usedTermQty := 0.0
	for _, order := range dispatcher.newOrders {
		usedTermQty += order.Qty * order.Price
	}
	assert.LessOrEqual(t, usedTermQty, termInventory)

	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	// no sleep on purpose

	// Move mid back down to the previous mid after not sleeping to hit the fire wall minimum quote time
	mid = 10_080.0

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	now = now.Add(5 * time.Millisecond)

	// Move the bid back down to previous mid and confirm before adding layers
	mid = 10_080.0

	// Run bids for reset
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)
	assert.Equal(t, 3, len(dispatcher.newOrders))
	assert.Equal(t, 3, len(dispatcher.cancelOrders))

	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Add another layer that breaches the inventory
	baseQtyLevels = append(baseQtyLevels, 5.0)
	model.SetLevelQtysInBase(baseQtyLevels)

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders)) // no cancel because we are adding a new layer
	assert.Equal(t, 10_291.68, dispatcher.newOrders[0].Price)
	assert.Equal(t, 4.0, dispatcher.newOrders[0].Qty) // 4.0 because we are breaching the inventory
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders)) // no cancel because we are adding a new layer
	assert.Equal(t, 9_868.32, dispatcher.newOrders[0].Price)
	assert.Equal(t, 4.118, dispatcher.newOrders[0].Qty) // 4.118 because we are breaching the inventory
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	usedTermQty = 0.00
	iter := orderBook.Bids.Iterator()
	for iter.Next() {
		usedTermQty += iter.Value().Price * iter.Value().Qty
	}
	assert.LessOrEqual(t, usedTermQty, termInventory)
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Move the mid up so we have to move the first layer and adjust the outer layer that is breaching the inventory
	mid = 10_100.0

	// Skip the Offers

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 2, len(dispatcher.newOrders))
	assert.Equal(t, 2, len(dispatcher.cancelOrders))
	assert.Equal(t, 10_089.90, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 9_887.90, dispatcher.newOrders[1].Price) // replace the outer layer since we were breaching the inventory
	assert.Equal(t, 4.108, dispatcher.newOrders[1].Qty)      // 4.108 because we are breaching the inventory
	assert.Equal(t, 10_069.92, dispatcher.cancelOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.cancelOrders[0].Qty)
	// Canceled outer layer
	assert.Equal(t, 9_868.32, dispatcher.cancelOrders[1].Price)
	assert.Equal(t, 4.118, dispatcher.cancelOrders[1].Qty)

	iter = orderBook.Bids.Iterator()
	usedTermQty = 0.00
	for iter.Next() {
		usedTermQty += iter.Value().Price * iter.Value().Qty
	}
	assert.LessOrEqual(t, usedTermQty, termInventory)
	dispatcher.Reset()

	// Cancel all orders
	makerEngine.CancelAll()
	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 8, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}

	bidIter := orderBook.Bids.Iterator()
	offerIter := orderBook.Offers.Iterator()
	assert.False(t, bidIter.Next())
	assert.False(t, offerIter.Next())

	baseQtyLevels = []float64{1.0, 0.25, 5.0}
	model.SetLevelQtysInBase(baseQtyLevels)

	mid = 10_000.0
	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)
	assert.Equal(t, 3, len(dispatcher.newOrders))

	assert.Equal(t, 10_010.0, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_012.5, dispatcher.newOrders[1].Price)
	assert.Equal(t, 0.25, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 10_062.5, dispatcher.newOrders[2].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[2].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	mid = 10_003.34
	// Run Offers again - this should force us to replace the next layer up as well
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)
	assert.Equal(t, 2, len(dispatcher.newOrders))

	assert.Equal(t, 10_013.35, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 10_015.85, dispatcher.newOrders[1].Price)
	assert.Equal(t, 0.25, dispatcher.newOrders[1].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	baseQtyLevels = []float64{1.0, 0.25, 0.1}
	model.SetLevelQtysInBase(baseQtyLevels)

	// Run Offers to replace the last layer
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)
	assert.Equal(t, 1, len(dispatcher.newOrders))

	assert.Equal(t, 10_016.85, dispatcher.newOrders[0].Price)
	assert.Equal(t, 0.1, dispatcher.newOrders[0].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	mid = 10_006.675

	// Run Offers to see all layers replaced despite not moving by 1/3 spread for all layers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)
	assert.Equal(t, 3, len(dispatcher.newOrders))
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	// Now do bids

	mid = 10_000.0
	baseQtyLevels = []float64{1.0, 0.25, 5.0}
	model.SetLevelQtysInBase(baseQtyLevels)

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)
	assert.Equal(t, 3, len(dispatcher.newOrders))

	assert.Equal(t, 9_990.0, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 9_987.5, dispatcher.newOrders[1].Price)
	assert.Equal(t, 0.25, dispatcher.newOrders[1].Qty)
	assert.Equal(t, 9_937.5, dispatcher.newOrders[2].Price)
	assert.Equal(t, 5.0, dispatcher.newOrders[2].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	mid = 9_996.66

	// Run Bids again - this should force us to replace the next layer up as well
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)
	assert.Equal(t, 2, len(dispatcher.newOrders))

	assert.Equal(t, 9_986.66, dispatcher.newOrders[0].Price)
	assert.Equal(t, 1.0, dispatcher.newOrders[0].Qty)
	assert.Equal(t, 9_984.16, dispatcher.newOrders[1].Price)
	assert.Equal(t, 0.25, dispatcher.newOrders[1].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	baseQtyLevels = []float64{1.0, 0.25, 0.1}
	model.SetLevelQtysInBase(baseQtyLevels)

	// Run Bids to replace the last layer
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)
	assert.Equal(t, 1, len(dispatcher.newOrders))

	assert.Equal(t, 9_983.16, dispatcher.newOrders[0].Price)
	assert.Equal(t, 0.1, dispatcher.newOrders[0].Qty)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

	mid = 9_993.32

	// Run Bids to see all layers replaced despite not moving by 1/3 spread for all layers
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)
	assert.Equal(t, 3, len(dispatcher.newOrders))
	dispatcher.Reset()
	now = now.Add(5 * time.Millisecond)

}

func TestMakerEngine_Replenishment(t *testing.T) {
	spot := instr.Spot{Base: instr.BTC, Term: instr.USD}
	ex := exchange.Coinbase
	pqRefData := &refdata.Precision{
		PricePrecision: 2,
		QtyPrecision:   3,
	}
	tradeSizeLimit := &refdata.BaseTradeSizeLimit{
		BaseMinTradeSize:      0.001,
		BaseMaxTradeSize:      1000,
	}
	refData := refdata.Composite{
		PriceQty: pqRefData,
		TradeSizeLimit: tradeSizeLimit,
		MinQuoteTime:          1 * time.Millisecond,
		ReplenishmentRate:     10 * time.Millisecond,
		UsePostOnly:          false,
	}

	mid := 10_000.0
	baseInventory := 20.0
	termInventory := 200_000.0
	baseQtyLevels := []float64{1, 3, 5}
	now := time.Now()

	orderBook := maker.NewOrderBook(spot, ex)
	dispatcher := NewMockDispatcher[instr.Spot]()
	controller := maker.NewMakerLayerController(spot, ex)
	model := maker.NewFunctionalLiqudityCurve(
		maker.NewStaticSymmetricCurveModel(
			func(totalBaseQty float64) float64 {
				return totalBaseQty * 10.0 // 10 bps per 1 BTC
			}),
		spot,
		ex,
	)

	model.SetLevelQtysInBase(baseQtyLevels)
	makerEngine := maker.NewEngine[instr.Spot](dispatcher, orderBook, controller, model)

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 6, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 3, orderBook.Offers.Size())
	assert.Equal(t, 3, orderBook.Bids.Size())

	mid = 20_000.0 // move this drastically to show that we will hit the minimum quote time

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	// Run bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))

	now = now.Add(5 * time.Millisecond)

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 6, len(dispatcher.newOrders))
	assert.Equal(t, 6, len(dispatcher.cancelOrders))

	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 3, orderBook.Offers.Size())
	assert.Equal(t, 3, orderBook.Bids.Size())

	now = now.Add(5 * time.Millisecond)

	// Fill the first two bids
	iter := orderBook.Bids.Iterator()
	iter.Next()
	bid := iter.Value()
	assert.NotNil(t, bid)
	assert.Equal(t, 1.0, bid.Qty)
	fillReport1 := &report.ExecutionReport[instr.Spot]{
		Instrument:    spot,
		Exchange:      ex,
		Side:          bid.Side,
		OrderID:       bid.ClientOrderID,
		ClientOrderID: bid.ClientOrderID,
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Price:         bid.Price,
		Qty:           bid.Qty,
		TimeInForce:   bid.TimeInForce,
		TransactTime:  now,
	}

	iter.Next()
	bid = iter.Value()
	assert.NotNil(t, bid)
	assert.Equal(t, 3.0, bid.Qty)
	fillReport2 := &report.ExecutionReport[instr.Spot]{
		Instrument:    spot,
		Exchange:      ex,
		Side:          bid.Side,
		OrderID:       bid.ClientOrderID,
		ClientOrderID: bid.ClientOrderID,
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Price:         bid.Price,
		Qty:           bid.Qty,
		TimeInForce:   bid.TimeInForce,
		TransactTime:  now,
	}

	orderBook.ApplyExecutionReport(fillReport1)
	controller.ApplyExecutionReport(fillReport1)
	orderBook.ApplyExecutionReport(fillReport2)
	controller.ApplyExecutionReport(fillReport2)

	assert.Equal(t, 1, orderBook.Bids.Size())
	assert.Equal(t, 3, orderBook.Offers.Size())

	// Run Bids without moving the price immediately to show we won't replenish the bid until later
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	// Do the same after 5 milliseconds but have the market move which will replace the outer layer
	now = now.Add(5 * time.Millisecond)

	// Run Bids
	mid = mid * 1.20
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 1, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 1, orderBook.Bids.Size())
	assert.Equal(t, 3, orderBook.Offers.Size())

	// Do the same after another 5 milliseconds to show the replenishment
	now = now.Add(5 * time.Millisecond)

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 2, orderBook.Bids.Size())
	assert.Equal(t, 3, orderBook.Offers.Size())

	// Wait until the next replenishment, which breaches inventory so we also cancel and replace the outermost layer
	now = now.Add(10 * time.Millisecond)

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 2, len(dispatcher.newOrders))
	assert.Equal(t, 1, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 3, orderBook.Bids.Size())
	assert.Equal(t, 3, orderBook.Offers.Size())

	// move mid back to original price
	mid = 20_000.0

	// Run Offers to show no change
	makerEngine.ReconcileAndReplaceOffers(mid, termInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	// Fill the first two offers
	iter = orderBook.Offers.Iterator()
	iter.Next()
	offer := iter.Value()
	assert.NotNil(t, offer)
	assert.Equal(t, 1.0, offer.Qty)
	fillReport1 = &report.ExecutionReport[instr.Spot]{
		Instrument:    spot,
		Exchange:      ex,
		Side:          offer.Side,
		OrderID:       offer.ClientOrderID,
		ClientOrderID: offer.ClientOrderID,
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Price:         offer.Price,
		Qty:           offer.Qty,
		TimeInForce:   offer.TimeInForce,
		TransactTime:  now,
	}

	iter.Next()
	offer = iter.Value()
	assert.NotNil(t, offer)
	assert.Equal(t, 3.0, offer.Qty)
	fillReport2 = &report.ExecutionReport[instr.Spot]{
		Instrument:    spot,
		Exchange:      ex,
		Side:          offer.Side,
		OrderID:       offer.ClientOrderID,
		ClientOrderID: offer.ClientOrderID,
		OrderStatus:   ordstatus.Filled,
		OrderType:     ordtypes.Limit,
		Price:         offer.Price,
		Qty:           offer.Qty,
		TimeInForce:   offer.TimeInForce,
		TransactTime:  now,
	}

	orderBook.ApplyExecutionReport(fillReport1)
	controller.ApplyExecutionReport(fillReport1)
	orderBook.ApplyExecutionReport(fillReport2)
	controller.ApplyExecutionReport(fillReport2)

	assert.Equal(t, 3, orderBook.Bids.Size())
	assert.Equal(t, 1, orderBook.Offers.Size())

	// Run Offers instantly to show no change
	makerEngine.ReconcileAndReplaceOffers(mid, termInventory, &refData, now)

	assert.Equal(t, 0, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	dispatcher.Reset()

	// Move the mid drastically to replace the outermost offer, but won't replenish due to the time constraint
	now = now.Add(5 * time.Millisecond)
	mid = mid * 1.5

	// Run Offers to replace the outermost offer
	makerEngine.ReconcileAndReplaceOffers(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 1, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 3, orderBook.Bids.Size())
	assert.Equal(t, 1, orderBook.Offers.Size())

	// Sleep to allow the replenish to happen
	now = now.Add(5 * time.Millisecond)

	// Run Offers to show replenishment
	makerEngine.ReconcileAndReplaceOffers(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 3, orderBook.Bids.Size())
	assert.Equal(t, 2, orderBook.Offers.Size())

	// Sleep for the total replenishment rate period to show the top of book offer being placed
	now = now.Add(10 * time.Millisecond)

	// Run Offers
	makerEngine.ReconcileAndReplaceOffers(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
	for _, order := range dispatcher.cancelOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.Canceled,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	for _, order := range dispatcher.newOrders {
		execReport := &report.ExecutionReport[instr.Spot]{
			Instrument:    spot,
			Exchange:      ex,
			Side:          order.Side,
			OrderID:       order.ClientOrderID,
			ClientOrderID: order.ClientOrderID,
			OrderStatus:   ordstatus.New,
			OrderType:     ordtypes.Limit,
			Price:         order.Price,
			Qty:           order.Qty,
			TimeInForce:   order.TimeInForce,
			TransactTime:  now,
		}
		orderBook.ApplyExecutionReport(execReport)
		controller.ApplyExecutionReport(execReport)
	}
	dispatcher.Reset()

	assert.Equal(t, 3, orderBook.Bids.Size())
	assert.Equal(t, 3, orderBook.Offers.Size())

}

func TestMakerEngine_USDT_USDC_Flickering(t *testing.T) {
	spot := instr.Spot{Base: instr.USDT, Term: instr.USDC}
	src := source.Coinbase
	ex := exchange.Coinbase
	pqRefData := &refdata.Precision{
		PricePrecision: 4,
		QtyPrecision:   2,
	}
	tradeSizeLimit := &refdata.BaseTradeSizeLimit{
		BaseMinTradeSize:      0.0001,
		BaseMaxTradeSize:      1_000_000.0,
	}
	refData := refdata.Composite{
		PriceQty: pqRefData,
		TradeSizeLimit: tradeSizeLimit,
		MinQuoteTime:          1 * time.Millisecond,
		ReplenishmentRate:     10 * time.Millisecond,
		UsePostOnly:          false,
	}

	tick := ticker.Ticker[instr.Spot]{
		Instrument: spot,
		Source:   src,
		Bid:        1.00059,
		Offer:      1.00060,
		Last:       1.00060,
	}
	baseInventory := 19_995.0
	termInventory := 19_995.0
	baseQtyLevels := []float64{20_000.0}
	now := time.Now()

	orderBook := maker.NewOrderBook(spot, ex)
	dispatcher := NewMockDispatcher[instr.Spot]()
	controller := maker.NewMakerLayerController(spot, ex)

	spreadModel := maker.NewWithBiasSpreadModel(maker.NewTwoFactorSymmetricCurveModel())
	liquidityCurve := maker.NewFunctionalLiqudityCurve(spreadModel, spot, exchange.Dexalot)

	makerEngine := maker.NewEngine[instr.Spot](dispatcher, orderBook, controller, liquidityCurve)
	spreadModel.SetBidBiasBps(-1.0)
	spreadModel.SetOfferBiasBps(-1.0)
    spreadModel.UnderlyingModel.SetInsideSpreadBps(1.0)
	spreadModel.UnderlyingModel.SetInsideQty(20_000.0)
	spreadModel.UnderlyingModel.SetQtySpreadMultiplierBps(1.0)
	liquidityCurve.SetLevelQtysInBase(baseQtyLevels)

	// Get the mid
	mid := price.TwoWayFromTicker(tick).Mid()
    assert.InDelta(t, 1.000595, mid, 0.0000001)
    mid = price.Skew(mid, -1.0)
    assert.InDelta(t, 1.00049, mid, 0.00001)

	// Run Bids
	makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
    assert.Equal(t, 1.00040, dispatcher.newOrders[0].Price)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
	dispatcher.Reset()

	// Run Offers
    now = now.Add(10 * time.Millisecond)
	makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

	assert.Equal(t, 1, len(dispatcher.newOrders))
	assert.Equal(t, 0, len(dispatcher.cancelOrders))
    assert.Equal(t, 1.00050, dispatcher.newOrders[0].Price)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
    dispatcher.Reset()

    // Run Offers 
    makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

    assert.Equal(t, 0, len(dispatcher.newOrders))
    assert.Equal(t, 0, len(dispatcher.cancelOrders))

    tick = ticker.Ticker[instr.Spot]{
        Instrument: spot,
        Source:   src,
        Bid:        1.00058,
        Offer:      1.00060,
        Last:       1.00051,
    }

    // Get the mid 
    mid = price.TwoWayFromTicker(tick).Mid()
    assert.InDelta(t, 1.00059, mid, 0.0000001)
    mid = price.Skew(mid, -1.0)
    assert.InDelta(t, 1.00049, mid, 0.00001)
    now = now.Add(10 * time.Millisecond)
    
    // Run Bids 
    makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now) 

    assert.Equal(t, 0, len(dispatcher.newOrders))
    assert.Equal(t, 0, len(dispatcher.cancelOrders))
    dispatcher.Reset()

    // Run Offers 
    makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

    assert.Equal(t, 0, len(dispatcher.newOrders))
    assert.Equal(t, 0, len(dispatcher.cancelOrders))
    dispatcher.Reset()

    tick = ticker.Ticker[instr.Spot]{
        Instrument: spot,
        Source:   src,
        Bid:        1.00060,
        Offer:      1.00061,
        Last:       1.00061,
    }

    // Get the mid 
    mid = price.TwoWayFromTicker(tick).Mid()
    assert.InDelta(t, 1.000605, mid, 0.0000001)
    mid = price.Skew(mid, -1.0)
    assert.InDelta(t, 1.000495, mid, 0.00001)
    now = now.Add(10 * time.Millisecond)

    // Run Bids 
    makerEngine.ReconcileAndReplaceBids(mid, termInventory, &refData, now)

    assert.Equal(t, 1, len(dispatcher.newOrders))
    assert.Equal(t, 1, len(dispatcher.cancelOrders))
    assert.Equal(t, 1.00050, dispatcher.newOrders[0].Price)
	for _, order := range dispatcher.cancelOrders {
		orderBook.Bids.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
	}
	for _, order := range dispatcher.newOrders {
		orderBook.Bids.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

	}
    dispatcher.Reset()

    // Run Offers 
    makerEngine.ReconcileAndReplaceOffers(mid, baseInventory, &refData, now)

    assert.Equal(t, 1, len(dispatcher.newOrders))
    assert.Equal(t, 1, len(dispatcher.cancelOrders))
    assert.Equal(t, 1.00060, dispatcher.newOrders[0].Price)
    for _, order := range dispatcher.cancelOrders {
        orderBook.Offers.Remove(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID})
    }
    for _, order := range dispatcher.newOrders {
        orderBook.Offers.Put(maker.Layer{Price: order.Price, ClientOrderID: order.ClientOrderID}, order)

    }
    dispatcher.Reset()
}
