package floats_test

import (
	"math"
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/floats"
)

func TestRoundWithPrecision(t *testing.T) {
	val := 1.245
	prec := 2
	rounded := floats.RoundWithPrecision(val, prec)
	if rounded != 1.25 {
		t.Errorf("expected %f, got %f", 1.25, rounded)
	}

	val = -.13765
	prec = 3
	rounded = floats.RoundWithPrecision(val, prec)
	if rounded != -0.138 {
		t.Errorf("expected %f, got %f", -0.138, rounded)
	}

	val = 1246782.1
	prec = 0
	rounded = floats.RoundWithPrecision(val, prec)
	if rounded != 1246782 {
		t.Errorf("expected %f, got %f", 1246782.0, rounded)
	}

	val = 1246782.9
	prec = 0
	rounded = floats.RoundWithPrecision(val, prec)
	if rounded != 1246783 {
		t.Errorf("expected %f, got %f", 1246783.0, rounded)
	}

	val = 54.2345
	prec = 6
	rounded = floats.RoundWithPrecision(val, prec)
	if rounded != 54.2345 {
		t.Errorf("expected %f, got %f", 54.2345, rounded)
	}

	val = 0.0
	prec = 2
	rounded = floats.RoundWithPrecision(val, prec)
	if rounded != 0.0 {
		t.Errorf("expected %f, got %f", 0.0, rounded)
	}

	val = -10.4
	prec = 0
	rounded = floats.RoundWithPrecision(val, prec)
	if rounded != -10 {
		t.Errorf("expected %f, got %f", -10.0, rounded)
	}

	val = -10.5
	prec = 0
	rounded = floats.RoundWithPrecision(val, prec)
	if rounded != -11 {
		t.Errorf("expected %f, got %f", -10.0, rounded)
	}
}

func TestRoundDirectionalWithPrecision(t *testing.T) {
	val := 1.245
	prec := 2
	rounded := floats.RoundDirectionalWithPrecision(val, prec, floats.RoundUp)
	if rounded != 1.25 {
		t.Errorf("expected %f, got %f", 1.25, rounded)
	}

	val = -.13765
	prec = 3
	rounded = floats.RoundDirectionalWithPrecision(val, prec, floats.RoundDown)
	if rounded != -0.138 {
		t.Errorf("expected %f, got %f", -0.138, rounded)
	}

	val = 1246782.1
	prec = 0
	rounded = floats.RoundDirectionalWithPrecision(val, prec, floats.RoundUp)
	if rounded != 1246783 {
		t.Errorf("expected %f, got %f", 1246783.0, rounded)
	}

	val = 1246782.9
	prec = 0
	rounded = floats.RoundDirectionalWithPrecision(val, prec, floats.RoundDown)
	if rounded != 1246782 {
		t.Errorf("expected %f, got %f", 1246782.0, rounded)
	}

	val = 54.2345
	prec = 6
	rounded = floats.RoundDirectionalWithPrecision(val, prec, floats.RoundUp)
	if rounded != 54.2345 {
		t.Errorf("expected %f, got %f", 54.2345, rounded)
	}

	val = 0.0
	prec = 2
	rounded = floats.RoundDirectionalWithPrecision(val, prec, floats.RoundUp)
	if rounded != 0.0 {
		t.Errorf("expected %f, got %f", 0.0, rounded)
	}

	val = -10.4
	prec = 0
	rounded = floats.RoundDirectionalWithPrecision(val, prec, floats.RoundUp)
	if rounded != -10 {
		t.Errorf("expected %f, got %f", -10.0, rounded)
	}

	val = -10.5
	prec = 0
	rounded = floats.RoundDirectionalWithPrecision(val, prec, floats.RoundDown)
	if rounded != -11 {
		t.Errorf("expected %f, got %f", -11.0, rounded)
	}
}

func TestRoundWithIncrement(t *testing.T) {
	cases := []struct {
		value     float64
		increment float64
		expected  float64
	}{
		{1.234, 0.1, 1.2},
		{1.234, 0.01, 1.23},
		{1.236, 0.01, 1.24},
		{1.234, 0.05, 1.25},
		{1.234, 0.025, 1.225},
		{1.47, 0.03125, 1.46875},
		{1.236, 0.001, 1.236},
		{1.234, 0.00025, 1.234},
	}

	for _, c := range cases {
		got := floats.RoundWithIncrement(c.value, c.increment)
		if math.Abs(got-c.expected) > 0.00000001 {
			t.Errorf("RoundFloat64Increment(%f, %f) = %f; want %f", c.value, c.increment, got, c.expected)
		}
	}
}

func TestRoundDirectionalWithIncrement(t *testing.T) {
	cases := []struct {
		value     float64
		increment float64
		roundMode floats.Mode
		expected  float64
	}{
		{1.234, 0.1, floats.RoundDown, 1.2},
		{1.234, 0.1, floats.RoundUp, 1.3},
		{1.234, 0.01, floats.RoundDown, 1.23},
		{1.234, 0.01, floats.RoundUp, 1.24},
		{1.234, 0.05, floats.RoundDown, 1.2},
		{1.234, 0.05, floats.RoundUp, 1.25},
		{1.234, 0.025, floats.RoundDown, 1.225},
		{1.234, 0.025, floats.RoundUp, 1.25},
		{1.47, 0.03125, floats.RoundDown, 1.46875},
		{1.47, 0.03125, floats.RoundUp, 1.5},
		{1.236, 0.00025, floats.RoundDown, 1.236},
		{1.236, 0.00025, floats.RoundUp, 1.236},
		{1.2354, 0.00025, floats.RoundDown, 1.23525},
		{1.2354, 0.00025, floats.RoundUp, 1.2355},
	}

	for _, c := range cases {
		got := floats.RoundDirectionalWithIncrement(c.value, c.increment, c.roundMode)
		if math.Abs(got-c.expected) > 0.00000001 {
			t.Errorf("DirectionalRoundFloat64Increment(%f, %f, %d) = %f; want %f", c.value, c.increment, c.roundMode, got, c.expected)
		}
	}
}

func TestIncrementDecimals(t *testing.T) {
	tests := []struct {
		name      string
		increment float64
		want      int
	}{
		{
			name:      "one decimal",
			increment: 0.1,
			want:      1,
		},
		{
			name:      "two decimals",
			increment: 0.01,
			want:      2,
		},
		{
			name:      "three decimals",
			increment: 0.001,
			want:      3,
		},
		{
			name:      "non power of 10 increment",
			increment: 0.00025,
			want:      5,
		},
		{
			name:      "non power of 10 increment",
			increment: 0.000250000,
			want:      5,
		},
		{
			name:      "non power of 10 increment",
			increment: 0.000250001,
			want:      9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := floats.IncrementToPrecision(tt.increment); got != tt.want {
				t.Errorf("Decimals() = %v, want %v", got, tt.want)
			}
		})
	}
}
