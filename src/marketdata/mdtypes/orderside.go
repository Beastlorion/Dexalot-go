package mdtypes

import "errors"

type OrderSide string

const (
	NilSide OrderSide = ""
	BUY     OrderSide = "BUY"
	SELL    OrderSide = "SELL"
)

var SideStringError = errors.New("invalid order side")

func SideFromString(s string) (OrderSide, error) {
	if s == "BUY" {
		return BUY, nil
	} else if s == "SELL" {
		return SELL, nil
	} else {
		return NilSide, SideStringError
	}
}

func (side OrderSide) Reverse() OrderSide {
	switch side {
	case BUY:
		return SELL
	case SELL:
		return BUY
	default:
		return NilSide
	}
}

func (side OrderSide) Sign() int {
	switch side {
	case BUY:
		return 1
	case SELL:
		return -1
	default:
		return 0
	}
}

