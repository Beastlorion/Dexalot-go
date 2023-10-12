package price_test

import (
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/price"
)

func TestEquals(t *testing.T) {
	price1 := price.TwoWay[instr.Spot]{
		Instrument:   instr.Spot{Base: instr.AVAX, Term: instr.USD},
		Bid:          19.0,
		Offer:        21.0,
		BidQty:       100.0,
		OfferQty:     125.0,
		TransactTime: time.UnixMilli(10_000_000),
	}

	price2 := price.TwoWay[instr.Spot]{
		Instrument:   instr.Spot{Base: instr.AVAX, Term: instr.USD},
		Bid:          19.0,
		Offer:        21.0,
		BidQty:       100.0,
		OfferQty:     125.0,
		TransactTime: time.UnixMilli(10_000_000),
	}

	price3 := price.TwoWay[instr.Spot]{
		Instrument:   instr.Spot{Base: instr.AVAX, Term: instr.USD},
		Bid:          20.0,
		Offer:        21.0,
		BidQty:       100.0,
		OfferQty:     125.0,
		TransactTime: time.UnixMilli(10_000_000),
	}

	price4 := price.TwoWay[instr.Spot]{}

	if !price1.Equals(price1) {
		t.Errorf("Equals Test FAILED. Same TwoWays were not equal")
	}

	if !price1.Equals(price2) {
		t.Errorf("Equals Test FAILED. Equal TwoWays were not equal")
	}

	if !price2.Equals(price1) {
		t.Errorf("Equals Test FAILED. Equal TwoWays were not equal")
	}

	if price1.Equals(price3) {
		t.Errorf("Equals Test FAILED. Unequal TwoWays were equal")
	}

	if price3.Equals(price1) {
		t.Errorf("Equals Test FAILED. Unequal TwoWays were equal")
	}

	if price1.Equals(price4) {
		t.Errorf("Equals Test FAILED. Empty price was equal")
	}
}

func TestMid(t *testing.T) {
	price := price.TwoWay[instr.Spot]{
		Instrument:   instr.Spot{Base: instr.AVAX, Term: instr.USD},
		Bid:          19.0,
		Offer:        21.0,
		BidQty:       100.0,
		OfferQty:     125.0,
		TransactTime: time.UnixMilli(10_000_000),
	}

	expected := 20.0
	if expected != price.Mid() {
		t.Errorf("Mid Test FAILED. Expected %f, got %f", expected, price.Mid())
	}
}

func TestRawSpread(t *testing.T) {
	price := price.TwoWay[instr.Spot]{
		Instrument:   instr.Spot{Base: instr.AVAX, Term: instr.USD},
		Bid:          19.0,
		Offer:        21.0,
		BidQty:       100.0,
		OfferQty:     125.0,
		TransactTime: time.UnixMilli(10_000_000),
	}

	expected := 2.0
	if expected != price.RawSpread() {
		t.Errorf("Raw Spread Test FAILED. Expected %f, got %f", expected, price.RawSpread())
	}
}

func TestSpreadBps(t *testing.T) {
	price := price.TwoWay[instr.Spot]{
		Instrument:   instr.Spot{Base: instr.AVAX, Term: instr.USD},
		Bid:          19.0,
		Offer:        21.0,
		BidQty:       100.0,
		OfferQty:     125.0,
		TransactTime: time.UnixMilli(10_000_000),
	}

	expected := 1_000.0
	if expected != price.SpreadBps() {
		t.Errorf("Raw Spread Test FAILED. Expected %f, got %f", expected, price.SpreadBps())
	}
}
