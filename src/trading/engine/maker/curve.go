package maker

import (
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/price"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
)

type FunctionalLiquidityCurve[T instr.Instrument] struct {
	spreadModel     CumulativeQtySpreadModel
	instrument      T
	exchange        exchange.Exchange
	levelQtysInBase []float64
	modelBids       []*orders.Limit[T]
	modelOffers     []*orders.Limit[T]
}

func (lc *FunctionalLiquidityCurve[T]) GenerateBids(mid float64, now time.Time) []*orders.Limit[T] {
	lc.modelBids = lc.modelBids[:0]
	totalQtyInBase := 0.0
	for i := 0; i < len(lc.levelQtysInBase); i++ {
		totalQtyInBase += lc.levelQtysInBase[i]
		lc.modelBids = append(lc.modelBids, &orders.Limit[T]{
			Instrument:       lc.instrument,
			ExchangeMetaData: &exchange.MetaData{Exchange: lc.exchange},
			ClientOrderID:    orders.GenerateClientOrderID(),
			Price:            price.AddSpread(mid, lc.spreadModel.BidSpreadBps(totalQtyInBase), price.Down),
			Qty:              lc.levelQtysInBase[i],
			Side:             side.BUY,
			TransactTime:     now,
		})
	}
	return lc.modelBids
}

func (lc *FunctionalLiquidityCurve[T]) GenerateOffers(mid float64, now time.Time) []*orders.Limit[T] {
	lc.modelOffers = lc.modelOffers[:0]
	totalQtyInBase := 0.0
	for i := 0; i < len(lc.levelQtysInBase); i++ {
		totalQtyInBase += lc.levelQtysInBase[i]
		lc.modelOffers = append(lc.modelOffers, &orders.Limit[T]{
			Instrument:       lc.instrument,
			ExchangeMetaData: &exchange.MetaData{Exchange: lc.exchange},
			ClientOrderID:    orders.GenerateClientOrderID(),
			Price:            price.AddSpread(mid, lc.spreadModel.OfferSpreadBps(totalQtyInBase), price.Up),
			Qty:              lc.levelQtysInBase[i],
			Side:             side.SELL,
			TransactTime:     now,
		})
	}
	return lc.modelOffers
}

func (lc *FunctionalLiquidityCurve[T]) SetLevelQtysInBase(levelQtysInBase []float64) {
	lc.levelQtysInBase = levelQtysInBase
}

func NewFunctionalLiqudityCurve[T instr.Instrument](spreadModel CumulativeQtySpreadModel, instrument T, exchange exchange.Exchange) *FunctionalLiquidityCurve[T] {
	return &FunctionalLiquidityCurve[T]{
		spreadModel:     spreadModel,
		instrument:      instrument,
		exchange:        exchange,
		levelQtysInBase: make([]float64, 0),
		modelBids:       make([]*orders.Limit[T], 0),
		modelOffers:     make([]*orders.Limit[T], 0),
	}
}

type SingleLayerLiquidityCurve[T instr.Instrument] struct {
	staticSpreadModel CumulativeQtySpreadModel
	instrument        T
	exchange          exchange.Exchange
	modelBids         []*orders.Limit[T]
	modelOffers       []*orders.Limit[T]
	baseInventory     float64
	termInventory     float64
}

func (lc *SingleLayerLiquidityCurve[T]) GenerateBids(mid float64, now time.Time) []*orders.Limit[T] {
	lc.modelBids = lc.modelBids[:0]
	price := price.AddSpread(mid, lc.staticSpreadModel.BidSpreadBps(0), price.Down)
	lc.modelBids = append(lc.modelBids, &orders.Limit[T]{
		Instrument:       lc.instrument,
		ExchangeMetaData: &exchange.MetaData{Exchange: lc.exchange},
		ClientOrderID:    orders.GenerateClientOrderID(),
		Price:            price,
		Qty:              lc.termInventory / price,
		Side:             side.BUY,
		TransactTime:     now,
	})
	return lc.modelBids
}

func (lc *SingleLayerLiquidityCurve[T]) GenerateOffers(mid float64, now time.Time) []*orders.Limit[T] {
	lc.modelOffers = lc.modelOffers[:0]
	lc.modelOffers = append(lc.modelOffers, &orders.Limit[T]{
		Instrument:       lc.instrument,
		ExchangeMetaData: &exchange.MetaData{Exchange: lc.exchange},
		ClientOrderID:    orders.GenerateClientOrderID(),
		Price:            price.AddSpread(mid, lc.staticSpreadModel.OfferSpreadBps(0), price.Up),
		Qty:              lc.baseInventory,
		Side:             side.SELL,
		TransactTime:     now,
	})
	return lc.modelOffers
}

func (lc *SingleLayerLiquidityCurve[T]) SetBaseInventory(baseInventory float64) {
	lc.baseInventory = baseInventory
}

func (lc *SingleLayerLiquidityCurve[T]) SetTermInventory(termInventory float64) {
	lc.termInventory = termInventory
}

func NewSingleLayerLiquidityCurve[T instr.Instrument](
	staticSpreadModel CumulativeQtySpreadModel,
	instrument T,
	exchange exchange.Exchange) *SingleLayerLiquidityCurve[T] {
	return &SingleLayerLiquidityCurve[T]{
		staticSpreadModel: staticSpreadModel,
		instrument:        instrument,
		exchange:          exchange,
		modelBids:         make([]*orders.Limit[T], 0, 1),
		modelOffers:       make([]*orders.Limit[T], 0, 1),
	}
}
