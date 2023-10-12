package price

import (
	"math"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/mdtypes"
)

type TwoWay[T instr.Instrument] struct {
	Instrument   T       
	Bid          float64 
	Offer        float64 
	BidQty       float64 
	OfferQty     float64 
	TransactTime time.Time   
}

func (price *TwoWay[T]) Mid() float64 {
	return (price.Bid + price.Offer) / 2
}

func (price *TwoWay[T]) RawSpread() float64 {
	return price.Offer - price.Bid
}

func (price *TwoWay[T]) SpreadBps() float64 {
	return 10_000 * (price.Offer - price.Bid) / price.Mid()
}

func (price *TwoWay[T]) Equals(other TwoWay[T]) bool {
	return price.Offer == other.Offer &&
		price.Bid == other.Bid && price.OfferQty == other.OfferQty && price.BidQty == other.BidQty
}

func EmptyTwoWay[T instr.Instrument]() TwoWay[T] {
	return TwoWay[T]{
		Bid:          math.NaN(),
		Offer:        math.NaN(),
		BidQty:       0.0,
		OfferQty:     0.0,
		TransactTime: time.Unix(0, 0),
	}
}

func IsEmptyTwoWay[T instr.Instrument](price TwoWay[T]) bool {
	return math.IsNaN(price.Bid) && math.IsNaN(price.Offer) && price.BidQty == 0.0 && price.OfferQty == 0.0
}

func TwoWayFromTicker[T instr.Instrument](ticker mdtypes.Ticker[T]) TwoWay[T] {
	return TwoWay[T]{
		Instrument:   ticker.Instrument,
		Bid:          ticker.Bid,
		Offer:        ticker.Offer,
		BidQty:       0.0,
		OfferQty:     0.0,
		TransactTime: ticker.TransactTime,
	}
}
