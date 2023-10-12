package instr_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
)

func TestToSpotUnderscoreSeparated(t *testing.T) {
	tests := []struct {
		input          string
		expectedResult *instr.Spot
		expectError    bool
	}{
		{"avax_usdc", &instr.Spot{Base: instr.AVAX, Term: instr.USDC}, false},
		{"eth_btc", &instr.Spot{Base: instr.ETH, Term: instr.BTC}, false},
		{"invalid", nil, true},
	}

	for _, test := range tests {
		result, err := instr.ToSpotUnderscoreSeparated(test.input)
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input: %s", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input: %s, error: %v", test.input, err)
			}
			if result != *test.expectedResult {
				t.Errorf("Expected %v, got %v", test.expectedResult, result)
			}
		}
	}
}

func TestToSpotSlashSeparated(t *testing.T) {
	tests := []struct {
		input          string
		expectedResult *instr.Spot
		expectError    bool
	}{
		{"AVAX/USDC", &instr.Spot{Base: instr.AVAX, Term: instr.USDC}, false},
		{"ETH/BTC", &instr.Spot{Base: instr.ETH, Term: instr.BTC}, false},
		{"invalid", nil, true},
	}

	for _, test := range tests {
		result, err := instr.ToSpotSlashSeparated(test.input)
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input: %s", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input: %s, error: %v", test.input, err)
			}
			if result != *test.expectedResult {
				t.Errorf("Expected %v, got %v", test.expectedResult, result)
			}
		}
	}
}

func TestSpotToUpperCaseUnderscoreSeparated(t *testing.T) {
	pair := instr.Spot{Base: instr.BTC, Term: instr.USD}
	expected := "BTC_USD"

	result := instr.SpotToUpperCaseUnderscoreSeparated(pair)

	if result != expected {
		t.Errorf("SpotToUpperCaseUnderscoreSeparated(%v) returned %v, expected %v", pair, result, expected)
	}
}

func TestSpotToSlashSeparated(t *testing.T) {
	pair := instr.Spot{Base: instr.BTC, Term: instr.USD}
	expected := "BTC/USD"

	result := instr.SpotToUpperCaseSlashSeparated(pair)

	if result != expected {
		t.Errorf("SpotToUpperCaseSlashSeparated(%v) returned %v, expected %v", pair, result, expected)
	}
}

func TestSpotToLowerCaseSlashSeparated(t *testing.T) {
	pair := instr.Spot{Base: instr.BTC, Term: instr.USD}
	expected := "btc/usd"

	result := instr.SpotToLowerCaseSlashSeparated(pair)

	if result != expected {
		t.Errorf("SpotToLowerCaseSlashSeparated(%v) returned %v, expected %v", pair, result, expected)
	}
}

func TestSpotToLowerCaseUnderscoreSeparated(t *testing.T) {
	pair := instr.Spot{Base: instr.BTC, Term: instr.USD}
	expected := "btc_usd"

	result := instr.SpotToLowerCaseUnderscoreSeparated(pair)

	if result != expected {
		t.Errorf("SpotToLowerCaseUnderscoreSeparated(%v) returned %v, expected %v", pair, result, expected)
	}
}

func TestSpotToUpperCaseDashSeparated(t *testing.T) {
	pair := instr.Spot{Base: instr.BTC, Term: instr.USD}
	expected := "BTC-USD"

	result := instr.SpotToUpperCaseDashSeparated(pair)

	if result != expected {
		t.Errorf("SpotToUpperCaseDashSeparated(%v) returned %v, expected %v", pair, result, expected)
	}
}

func TestSpotToLowerCaseDashSeparated(t *testing.T) {
	pair := instr.Spot{Base: instr.BTC, Term: instr.USD}
	expected := "btc-usd"

	result := instr.SpotToLowerCaseDashSeparated(pair)

	if result != expected {
		t.Errorf("SpotToLowerCaseDashSeparated(%v) returned %v, expected %v", pair, result, expected)
	}
}
