package floats

import "math"

func CmpWithinDelta(a, b, delta float64) bool {
	return math.Abs(a-b) < delta
}

func CmpWithinDeltaOrEqual(a, b, delta float64) bool {
	return math.Abs(a-b) <= delta
}

func CmpDescendingOrder(a, b float64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func CmpAscendingOrder(a, b float64) int {
	if a > b {
		return -1
	} else if a < b {
		return 1
	} else {
		return 0
	}
}
