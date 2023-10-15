package maker

import (
	"math"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/floats"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/structs/tree"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
)

const halfSpreadRepostThreshold = 1.0 / 3.0

type LiquidityCurve[T instr.Instrument] interface {
	GenerateBids(mid float64, now time.Time) []*orders.Limit[T]
	GenerateOffers(mid float64, now time.Time) []*orders.Limit[T]
}

type Dispatcher[T instr.Instrument] interface {
	DispatchLimitOrder(order *orders.Limit[T])
	DispatchCancelOrder(order *orders.Limit[T])
}

type Engine[T instr.Instrument] struct {
	dispatcher         Dispatcher[T]
	confirmedOrderBook *OrderBook[T]
	controller         *MakerLayerController[T]
	liquidityCurve     LiquidityCurve[T]
}

func (e *Engine[T]) cancelRemaining(iter tree.NodeIterator[Layer, *orders.Limit[T]]) {
	for iter.Next() {
		e.dispatcher.DispatchCancelOrder(iter.Value())
	}
}

func (e *Engine[T]) CancelAll() {
	e.cancelRemaining(e.confirmedOrderBook.Bids.Iterator())
	e.cancelRemaining(e.confirmedOrderBook.Offers.Iterator())
}

func (e *Engine[T]) normalizeAndValidateOrder(newOrder *orders.Limit[T], remainingInventoryInDealt float64, refData *refdata.Composite) error {
	if refData.UsePostOnly {
		newOrder.TimeInForce = tif.PO
	} else {
		newOrder.TimeInForce = tif.GTC
	}

	err := orders.BoundQtyInDealt(newOrder, remainingInventoryInDealt)
	if err != nil {
		return err
	}

	err = orders.NormalizeMakerWithRefData(newOrder, refData)
	if err != nil {
		return err
	}

	err = orders.ValidateLimitWithRefData(newOrder, refData)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine[T]) shouldReplaceOrder(newOrder, oldOrder, lastOrder *orders.Limit[T], mid, remainingInventoryInDealt float64, refData *refdata.Composite, now time.Time) bool {
    repostThreshold := halfSpreadRepostThreshold * math.Abs(mid-newOrder.Price)
    switch newOrder.Side {
    case side.BUY:
        if remainingInventoryInDealt < oldOrder.Qty*oldOrder.Price {
            return true
        } else if lastOrder != nil && floats.CmpDescendingOrder(lastOrder.Price, oldOrder.Price) <= 0 {
            return true
        } else if !floats.CmpWithinDeltaOrEqual(oldOrder.Price, newOrder.Price, repostThreshold) {
            if e.controller.HasMetMinimumQuoteTime(oldOrder, refData.MinQuoteTime, now) {
                return true
            }
        }
    case side.SELL:
        if remainingInventoryInDealt < oldOrder.Qty {
            return true
        } else if lastOrder != nil && floats.CmpAscendingOrder(lastOrder.Price, oldOrder.Price) <= 0 {
            return true
        } else if !floats.CmpWithinDeltaOrEqual(oldOrder.Price, newOrder.Price, repostThreshold) {
            if e.controller.HasMetMinimumQuoteTime(oldOrder, refData.MinQuoteTime, now) {
                return true
            }
        }
    }
    return false
}

func (e *Engine[T]) ReconcileAndReplaceOffers(mid, remainingBaseInventory float64, refData *refdata.Composite, now time.Time) error {
	iter := e.confirmedOrderBook.Offers.Iterator()
	newOffers := e.liquidityCurve.GenerateOffers(mid, now)
    var lastConfirmedOrder *orders.Limit[T]
	for i := range newOffers {
		err := e.normalizeAndValidateOrder(newOffers[i], remainingBaseInventory, refData)
		if err != nil {
			e.cancelRemaining(iter) // cancel outer layers of the side we are working on
			switch err {
			case orders.NonPositiveQtyError:
				return nil // we are done looping
			case orders.NonPositivePriceError:
				return nil // we are done looping
			default:
				return err // this indicates a bug
			}
		}

		switch e.controller.GetLayerStatus(i, newOffers[i].Side, refData, now) {
		case NotReplenishable:
			continue
		case Replenishable:
			e.dispatcher.DispatchLimitOrder(newOffers[i])
			remainingBaseInventory -= newOffers[i].Qty
		case Reconcilable:
			if iter.Next() {
				if e.shouldReplaceOrder(newOffers[i], iter.Value(), lastConfirmedOrder, mid, remainingBaseInventory, refData, now) {
					e.dispatcher.DispatchCancelOrder(iter.Value())
					e.dispatcher.DispatchLimitOrder(newOffers[i])
					remainingBaseInventory -= newOffers[i].Qty
                    lastConfirmedOrder = newOffers[i]
				} else {
					remainingBaseInventory -= iter.Value().Qty
                    lastConfirmedOrder = iter.Value()
				}
			} else {
				e.dispatcher.DispatchLimitOrder(newOffers[i])
				remainingBaseInventory -= newOffers[i].Qty
                lastConfirmedOrder = newOffers[i]
			}
		}
	}
	if iter.Next() {
		e.cancelRemaining(iter)
	}
	return nil
}

func (e *Engine[T]) ReconcileAndReplaceBids(mid, remainingTermInventory float64, refData *refdata.Composite, now time.Time) error {
	iter := e.confirmedOrderBook.Bids.Iterator()
	newBids := e.liquidityCurve.GenerateBids(mid, now)
    var lastConfirmedOrder *orders.Limit[T]
	for i := range newBids {
		err := e.normalizeAndValidateOrder(newBids[i], remainingTermInventory, refData)
		if err != nil {
			e.cancelRemaining(iter) // cancel outer layers of the side we are working on
			switch err {
			case orders.NonPositiveQtyError:
				return nil // we are done looping
			case orders.NonPositivePriceError:
				return nil // we are done looping
			default:
				return err // this indicates a bug
			}
		}

		switch e.controller.GetLayerStatus(i, newBids[i].Side, refData, now) {
		case NotReplenishable:
			continue
		case Replenishable:
			e.dispatcher.DispatchLimitOrder(newBids[i])
			remainingTermInventory -= newBids[i].Qty * newBids[i].Price
		case Reconcilable:
			if iter.Next() {
				if e.shouldReplaceOrder(newBids[i], iter.Value(), lastConfirmedOrder, mid, remainingTermInventory, refData, now) {
					e.dispatcher.DispatchCancelOrder(iter.Value())
					e.dispatcher.DispatchLimitOrder(newBids[i])
					remainingTermInventory -= newBids[i].Qty * newBids[i].Price
                    lastConfirmedOrder = newBids[i]
				} else {
					remainingTermInventory -= iter.Value().Qty * iter.Value().Price
                    lastConfirmedOrder = iter.Value()
				}
			} else {
				e.dispatcher.DispatchLimitOrder(newBids[i])
				remainingTermInventory -= newBids[i].Qty * newBids[i].Price
                lastConfirmedOrder = newBids[i]
			}
		}
	}
	if iter.Next() {
		e.cancelRemaining(iter)
	}
	return nil
}

func NewEngine[T instr.Instrument](
	dispatcher Dispatcher[T],
	confirmedOrderBook *OrderBook[T],
	controller *MakerLayerController[T],
	liquidityCurve LiquidityCurve[T]) *Engine[T] {
	return &Engine[T]{
		dispatcher:         dispatcher,
		confirmedOrderBook: confirmedOrderBook,
		controller:         controller,
		liquidityCurve:     liquidityCurve,
	}
}
