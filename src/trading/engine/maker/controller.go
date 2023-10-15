package maker

import (
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordstatus"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordtypes"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/refdata"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/report"
)

type LayerStatus int

const (
	NotReplenishable LayerStatus = (iota + 1)
	Replenishable
	Reconcilable
)

type MakerLayerController[T instr.Instrument] struct {
	instrument              T
	exchange                exchange.Exchange
	bidReplenishmentIndex   int
	offerReplenishmentIndex int
	bidControlTime          time.Time
	offerControlTime        time.Time
	clordIDCache            map[string]struct{}
}

func (rm *MakerLayerController[T]) HasMetMinimumQuoteTime(orderLayer *orders.Limit[T], minimumQuoteTime time.Duration, now time.Time) bool {
	return now.Sub(orderLayer.TransactTime) >= minimumQuoteTime
}

func (rm *MakerLayerController[T]) GetLayerStatus(layerIndex int, s side.OrderSide, refData *refdata.Composite, now time.Time) LayerStatus {
	switch s {
	case side.BUY:
		if layerIndex > rm.bidReplenishmentIndex {
			return Reconcilable
		} else if layerIndex < rm.bidReplenishmentIndex {
			return NotReplenishable
		}

		if now.Sub(rm.bidControlTime) >= refData.ReplenishmentRate {
			rm.bidReplenishmentIndex = rm.bidReplenishmentIndex - 1
			rm.bidControlTime = now
			return Replenishable
		}
		return NotReplenishable
	case side.SELL:
		if layerIndex > rm.offerReplenishmentIndex {
			return Reconcilable
		} else if layerIndex < rm.offerReplenishmentIndex {
			return NotReplenishable
		}

		if now.Sub(rm.offerControlTime) >= refData.ReplenishmentRate {
			rm.offerReplenishmentIndex = rm.offerReplenishmentIndex - 1
			rm.offerControlTime = now
			return Replenishable
		}
		return NotReplenishable
	default:
		return NotReplenishable
	}
}

func (rm *MakerLayerController[T]) ApplyExecutionReport(rep *report.ExecutionReport[T]) {
	if rep.Instrument != rm.instrument || rep.Exchange != rm.exchange || rep.OrderType != ordtypes.Limit {
		return
	}

	if rep.TimeInForce != tif.GTC && rep.TimeInForce != tif.PO {
		return
	}

	switch rep.OrderStatus {
	case ordstatus.New:
		rm.clordIDCache[rep.ClientOrderID] = struct{}{}
	case ordstatus.Expired:
		fallthrough
	case ordstatus.Replaced:
		fallthrough
	case ordstatus.DoneForDay:
		fallthrough
	case ordstatus.Canceled:
		delete(rm.clordIDCache, rep.ClientOrderID)
	case ordstatus.Filled:
		_, ok := rm.clordIDCache[rep.ClientOrderID]
		if !ok {
			return
		}
		if rep.Side == side.BUY {
			rm.bidReplenishmentIndex += 1
			rm.bidControlTime = rep.TransactTime
		} else if rep.Side == side.SELL {
			rm.offerReplenishmentIndex += 1
			rm.offerControlTime = rep.TransactTime
		}
		delete(rm.clordIDCache, rep.ClientOrderID)
	}
}

func NewMakerLayerController[T instr.Instrument](instrument T, ex exchange.Exchange) *MakerLayerController[T] {
	return &MakerLayerController[T]{
		instrument:              instrument,
		exchange:                ex,
		bidReplenishmentIndex:   -1,
		offerReplenishmentIndex: -1,
		clordIDCache:            make(map[string]struct{}),
	}
}
