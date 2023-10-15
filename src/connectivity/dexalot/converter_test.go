package dexalot_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/connectivity/dexalot"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
)

func TestHashSpot(t *testing.T) {
	spot := instr.Spot{Base: instr.AVAX, Term: instr.USDC}
	hashed, err := dexalot.CreateSpotHash(spot)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected := "0x415641582f555344430000000000000000000000000000000000000000000000"
	if hashed.Hex() != expected {
		t.Errorf("expected %s, got %s", expected, hashed)
	}
}

func TestApiAsset(t *testing.T) {
	asset := instr.AVAX
	apiAsset, err := dexalot.ConvertToAPIAsset(asset)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected := "AVAX"
	if apiAsset != expected {
		t.Errorf("expected %s, got %s", expected, apiAsset)
	}

	asset = instr.USDC
	apiAsset, err = dexalot.ConvertToAPIAsset(asset)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected = "USDC"
	if apiAsset != expected {
		t.Errorf("expected %s, got %s", expected, apiAsset)
	}

	asset = instr.USDT
	apiAsset, err = dexalot.ConvertToAPIAsset(asset)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected = "USDt"
	if apiAsset != expected {
		t.Errorf("expected %s, got %s", expected, apiAsset)
	}

	asset = instr.ALOT
	apiAsset, err = dexalot.ConvertToAPIAsset(asset)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected = "ALOT"
	if apiAsset != expected {
		t.Errorf("expected %s, got %s", expected, apiAsset)
	}

	asset = instr.BTCb
	apiAsset, err = dexalot.ConvertToAPIAsset(asset)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected = "BTC.b"
	if apiAsset != expected {
		t.Errorf("expected %s, got %s", expected, apiAsset)
	}

	asset = instr.WETHe
	apiAsset, err = dexalot.ConvertToAPIAsset(asset)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected = "WETH.e"
	if apiAsset != expected {
		t.Errorf("expected %s, got %s", expected, apiAsset)
	}
}

func TestApiInstrument(t *testing.T) {
	spot := instr.Spot{Base: instr.AVAX, Term: instr.USDC}
	apiInstrument, err := dexalot.ConvertSpotToAPIInstrument(spot)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected := "AVAX/USDC"
	if apiInstrument != expected {
		t.Errorf("expected %s, got %s", expected, apiInstrument)
	}

	spot = instr.Spot{Base: instr.AVAX, Term: instr.USDT}
	apiInstrument, err = dexalot.ConvertSpotToAPIInstrument(spot)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected = "AVAX/USDt"
	if apiInstrument != expected {
		t.Errorf("expected %s, got %s", expected, apiInstrument)
	}

	spot = instr.Spot{Base: instr.BTCb, Term: instr.USDC}
	apiInstrument, err = dexalot.ConvertSpotToAPIInstrument(spot)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	expected = "BTC.b/USDC"
	if apiInstrument != expected {
		t.Errorf("expected %s, got %s", expected, apiInstrument)
	}
}
