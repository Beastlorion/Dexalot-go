package mdtypes_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/mdtypes"
)

func TestOrderSide(t *testing.T) {
	mySide := mdtypes.SELL
	result := mySide.Reverse()
	if result != mdtypes.BUY {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", mdtypes.BUY, result)
	}

	result, err := mdtypes.SideFromString("SELL")
	if err != nil || result != mdtypes.SELL {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", mdtypes.SELL, result)
	}

	result, err = mdtypes.SideFromString("BUY")
	if err != nil || result != mdtypes.BUY {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", mdtypes.BUY, result)
	}

	result, err = mdtypes.SideFromString("NOT A SIDE")
	if err == nil || result != "" {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", "", result)
	}

	if mySide.Reverse() != mdtypes.BUY {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", mdtypes.BUY, result)
	} else if mySide.Sign() != -1 {
		t.Errorf("Order side Test FAILED. Expected %d, got %d", -1, mySide.Sign())
	}

	mySide = mySide.Reverse()
	if mySide.Reverse() != mdtypes.SELL {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", mdtypes.SELL, result)
	} else if mySide.Sign() != 1 {
		t.Errorf("Order side Test FAILED. Expected %d, got %d", 1, mySide.Sign())
	}
}
