package floats

import (
	"math"
)

func UpperBound(value, upperBound float64) float64 {
	return math.Min(value, upperBound)
}

func LowerBound(value, lowerBound float64) float64 {
	return math.Max(value, lowerBound)
}

func UpperLowerBound(value, upperBound, lowerBound float64) float64 {
	return LowerBound(UpperBound(value, upperBound), lowerBound)
}
