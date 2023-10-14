package floats

import (
	"fmt"
	"math"
	"strings"
)

type Mode int

const (
	RoundUp   = 1
	RoundDown = -1
)

func RoundWithPrecision(value float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	rounded := math.Round(value*shift) / shift
	return rounded
}

func RoundDirectionalWithPrecision(f float64, precision int, roundingMode Mode) float64 {
	factor := math.Pow(10, float64(precision))
	if roundingMode == RoundUp {
		return math.Ceil(f*factor) / factor
	}
	return math.Floor(f*factor) / factor
}

func RoundWithIncrement(value float64, increment float64) float64 {
	if increment == 0.0 {
		return math.Round(value)
	}
	shift := 1.0 / increment
	return math.Round(value*shift) / shift
}

func RoundDirectionalWithIncrement(f float64, increment float64, roundingMode Mode) float64 {
	if increment == 0.0 {
		if roundingMode == RoundUp {
			return math.Ceil(f)
		} else {
			return math.Floor(f)
		}
	}
	shift := 1.0 / increment
	if roundingMode == RoundUp {
		return math.Ceil(f*shift) / shift
	}
	return math.Floor(f*shift) / shift
}

// It is best to do this once and store the result.
func IncrementToPrecision(increment float64) int {
	s := fmt.Sprintf("%.9f", increment)
	s = s[strings.Index(s, ".")+1:]
	s = strings.TrimRight(s, "0")
	return len(s)
}
