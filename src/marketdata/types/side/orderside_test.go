package side_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
)

func TestOrderSide(t *testing.T) {
	mySide := side.SELL
	result := mySide.Reverse()
	if result != side.BUY {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", side.BUY, result)
	}

	result, err := side.SideFromString("SELL")
	if err != nil || result != side.SELL {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", side.SELL, result)
	}

	result, err = side.SideFromString("BUY")
	if err != nil || result != side.BUY {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", side.BUY, result)
	}

	result, err = side.SideFromString("NOT A SIDE")
	if err == nil || result != "" {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", "", result)
	}

	if mySide.Reverse() != side.BUY {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", side.BUY, result)
	} else if mySide.Sign() != -1 {
		t.Errorf("Order side Test FAILED. Expected %d, got %d", -1, mySide.Sign())
	}

	mySide = mySide.Reverse()
	if mySide.Reverse() != side.SELL {
		t.Errorf("Order side Test FAILED. Expected %s, got %s", side.SELL, result)
	} else if mySide.Sign() != 1 {
		t.Errorf("Order side Test FAILED. Expected %d, got %d", 1, mySide.Sign())
	}
}
