package bytes32_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/bytes32"
)

func TestBytes32Equal(t *testing.T) {
	a := [32]byte{1, 2, 3, 4, 5}
	b := [32]byte{1, 2, 3, 4, 5}
	c := [32]byte{1, 2, 3, 4, 6}
	d := [32]byte{1, 2, 3, 4, 5, 7}
	e := [32]byte{1, 2, 3, 4}

	if !bytes32.Equal(a, b) {
		t.Error("Expected a and b to be equal")
	}

	if bytes32.Equal(a, c) {
		t.Error("Expected a and c to be different")
	}

	if bytes32.Equal(a, d) {
		t.Error("Expected a and d to be different")
	}

	if bytes32.Equal(b, c) {
		t.Error("Expected b and c to be different")
	}

	if bytes32.Equal(b, d) {
		t.Error("Expected b and d to be different")
	}

	if bytes32.Equal(c, d) {
		t.Error("Expected c and d to be different")
	}

	if bytes32.Equal(a, e) {
		t.Errorf("Expected a and e to be different")
	}

	if bytes32.Equal(b, e) {
		t.Errorf("Expected b and e to be different")
	}

	if bytes32.Equal(c, e) {
		t.Errorf("Expected c and e to be different")
	}

	if bytes32.Equal(d, e) {
		t.Errorf("Expected d and e to be different")
	}
}
