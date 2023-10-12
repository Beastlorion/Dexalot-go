package instr_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
)

func TestToInternalAssetName(t *testing.T) {
	tests := []struct {
		input          string
		expectedResult instr.Asset
	}{
		{"avax", instr.AVAX},
		{"UsDc", instr.USDC},
		{"USDC", instr.USDC},
		{"btc.b", instr.BTCb},
		{"btc.B", instr.BTCb},
		{"BTC.B", instr.BTCb},
		{"BTC.b", instr.BTCb},
		{"WETH.e", instr.WETHe},
		{"WETH.E", instr.WETHe},
		{"weth.e", instr.WETHe},
		{"weth.E", instr.WETHe},
	}

	for _, test := range tests {
		result := instr.ToInternalAssetName(test.input)
		if result != test.expectedResult {
			t.Errorf("Expected %s, got %s", test.expectedResult, result)
		}
	}
}
