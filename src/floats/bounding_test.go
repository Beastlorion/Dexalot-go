package floats_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/floats"
	"github.com/stretchr/testify/assert"
)

func TestUpperBoundFloat64(t *testing.T) {
	assert.Equal(t, 3.0, floats.UpperBound(3.0, 4.0))
	assert.Equal(t, 4.0, floats.UpperBound(4.0, 4.0))
	assert.Equal(t, 4.0, floats.UpperBound(5.0, 4.0))
}

func TestLowerBoundFloat64(t *testing.T) {
	assert.Equal(t, 3.0, floats.LowerBound(3.0, 2.0))
	assert.Equal(t, 2.0, floats.LowerBound(2.0, 2.0))
	assert.Equal(t, 2.0, floats.LowerBound(1.0, 2.0))
}

func TestUpperLowerBoundFloat64(t *testing.T) {
	assert.Equal(t, 3.0, floats.UpperLowerBound(3.0, 4.0, 2.0))
	assert.Equal(t, 4.0, floats.UpperLowerBound(5.0, 4.0, 2.0))
	assert.Equal(t, 2.0, floats.UpperLowerBound(1.0, 4.0, 2.0))
}
