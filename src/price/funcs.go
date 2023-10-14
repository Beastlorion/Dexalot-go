package price

import "math"

type SpreadDirection int

const bpsFactor = 10_000.0

const (
	Up = SpreadDirection(iota + 1)
	Down
)

func Skew(price, skewFactorBps float64) float64 {
	return price * (1 + skewFactorBps/bpsFactor)
}

func AddSpread(price, spreadBps float64, direction SpreadDirection) float64 {
	switch direction {
	case Up:
		return price * (1 + spreadBps/bpsFactor)
	case Down:
		return price * (1 - spreadBps/bpsFactor)
	default:
		panic("invalid direction")
	}
}

func GetSpreadBps(mid, price float64) float64 {
	return bpsFactor * (price - mid) / mid
}

func GetDistanceBps(mid, price1, price2 float64) float64 {
	return bpsFactor * (price2 - price1) / mid
}

func GetAbsoluteSpreadBps(mid, price float64) float64 {
	return math.Abs(GetSpreadBps(mid, price))
}

func GetAbsoluteDistanceBps(mid, price1, price2 float64) float64 {
	return math.Abs(GetDistanceBps(mid, price1, price2))
}
