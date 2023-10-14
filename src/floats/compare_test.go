package floats_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/floats"
	"github.com/stretchr/testify/assert"
)

func TestCmpWithinDelta(t *testing.T) {
	assert.False(t, floats.CmpWithinDelta(1.0, 2.0, 1.0))
	assert.False(t, floats.CmpWithinDelta(1.0, 2.0, 0.5))
	assert.True(t, floats.CmpWithinDelta(1.0, 2.0, 1.5))
	assert.True(t, floats.CmpWithinDelta(1.0, 2.0, 2.0))

	assert.False(t, floats.CmpWithinDelta(1.0, 1.0, 0.0))
	assert.True(t, floats.CmpWithinDelta(1.0, 1.0, 0.5))
	assert.True(t, floats.CmpWithinDelta(1.0, 1.0, 1.0))

	assert.False(t, floats.CmpWithinDelta(1.0, 0.0, 0.0))
	assert.False(t, floats.CmpWithinDelta(1.0, 0.0, 0.5))
	assert.True(t, floats.CmpWithinDelta(1.0, 0.0, 1.1))

	assert.False(t, floats.CmpWithinDelta(1.0, 1.0, -0.5))
	assert.False(t, floats.CmpWithinDelta(1.0, 0.5, -1.0))
}

func TestCmpWithinDeltaOrEqual(t *testing.T) {
	assert.True(t, floats.CmpWithinDeltaOrEqual(1.0, 2.0, 1.0))
	assert.False(t, floats.CmpWithinDeltaOrEqual(1.0, 2.0, 0.5))
	assert.True(t, floats.CmpWithinDeltaOrEqual(1.0, 2.0, 1.5))
	assert.True(t, floats.CmpWithinDeltaOrEqual(1.0, 2.0, 2.0))

	assert.True(t, floats.CmpWithinDeltaOrEqual(1.0, 1.0, 0.0))
	assert.True(t, floats.CmpWithinDeltaOrEqual(1.0, 1.0, 0.5))
	assert.True(t, floats.CmpWithinDeltaOrEqual(1.0, 1.0, 1.0))

	assert.False(t, floats.CmpWithinDeltaOrEqual(1.0, 0.0, 0.0))
	assert.False(t, floats.CmpWithinDeltaOrEqual(1.0, 0.0, 0.5))
	assert.True(t, floats.CmpWithinDeltaOrEqual(1.0, 0.0, 1.1))

	assert.False(t, floats.CmpWithinDeltaOrEqual(1.0, 1.0, -0.5))
	assert.False(t, floats.CmpWithinDeltaOrEqual(1.0, 0.5, -1.0))
}

func TestCmpCompareDescendingOrder(t *testing.T) {
	assert.True(t, floats.CmpDescendingOrder(1.0, 2.0) < 0)
	assert.False(t, floats.CmpDescendingOrder(2.0, 1.0) < 0)
	assert.True(t, floats.CmpDescendingOrder(2.0, 1.0) > 0)
	assert.False(t, floats.CmpDescendingOrder(1.0, 2.0) > 0)
	assert.True(t, floats.CmpDescendingOrder(1.0, 1.0) == 0)

	assert.True(t, floats.CmpDescendingOrder(1.0, 1.5) < 0)
	assert.False(t, floats.CmpDescendingOrder(1.5, 1.0) < 0)
	assert.True(t, floats.CmpDescendingOrder(1.5, 1.0) > 0)
	assert.False(t, floats.CmpDescendingOrder(1.0, 1.5) > 0)
	assert.True(t, floats.CmpDescendingOrder(1.5, 1.5) == 0)

	assert.True(t, floats.CmpDescendingOrder(1.0, 0.5) > 0)
	assert.False(t, floats.CmpDescendingOrder(0.5, 1.0) > 0)
	assert.True(t, floats.CmpDescendingOrder(0.5, 1.0) < 0)
	assert.False(t, floats.CmpDescendingOrder(1.0, 0.5) < 0)
	assert.True(t, floats.CmpDescendingOrder(0.5, 0.5) == 0)

	assert.True(t, floats.CmpDescendingOrder(1.0, 0.0) > 0)
	assert.False(t, floats.CmpDescendingOrder(0.0, 1.0) > 0)
	assert.True(t, floats.CmpDescendingOrder(0.0, 1.0) < 0)
	assert.False(t, floats.CmpDescendingOrder(1.0, 0.0) < 0)
	assert.True(t, floats.CmpDescendingOrder(0.0, 0.0) == 0)

	assert.True(t, floats.CmpDescendingOrder(1.0, -0.5) > 0)
	assert.False(t, floats.CmpDescendingOrder(-0.5, 1.0) > 0)
	assert.True(t, floats.CmpDescendingOrder(-0.5, 1.0) < 0)
	assert.False(t, floats.CmpDescendingOrder(1.0, -0.5) < 0)
	assert.True(t, floats.CmpDescendingOrder(-0.5, -0.5) == 0)
}

func TestCmpCompareAscendingOrder(t *testing.T) {
	assert.False(t, floats.CmpAscendingOrder(1.0, 2.0) < 0)
	assert.True(t, floats.CmpAscendingOrder(2.0, 1.0) < 0)
	assert.False(t, floats.CmpAscendingOrder(2.0, 1.0) > 0)
	assert.True(t, floats.CmpAscendingOrder(1.0, 2.0) > 0)
	assert.True(t, floats.CmpAscendingOrder(1.0, 1.0) == 0)

	assert.False(t, floats.CmpAscendingOrder(1.0, 1.5) < 0)
	assert.True(t, floats.CmpAscendingOrder(1.5, 1.0) < 0)
	assert.False(t, floats.CmpAscendingOrder(1.5, 1.0) > 0)
	assert.True(t, floats.CmpAscendingOrder(1.0, 1.5) > 0)
	assert.True(t, floats.CmpAscendingOrder(1.5, 1.5) == 0)

	assert.False(t, floats.CmpAscendingOrder(1.0, 0.5) > 0)
	assert.True(t, floats.CmpAscendingOrder(0.5, 1.0) > 0)
	assert.False(t, floats.CmpAscendingOrder(0.5, 1.0) < 0)
	assert.True(t, floats.CmpAscendingOrder(1.0, 0.5) < 0)
	assert.True(t, floats.CmpAscendingOrder(0.5, 0.5) == 0)

	assert.False(t, floats.CmpAscendingOrder(1.0, 0.0) > 0)
	assert.True(t, floats.CmpAscendingOrder(0.0, 1.0) > 0)
	assert.False(t, floats.CmpAscendingOrder(0.0, 1.0) < 0)
	assert.True(t, floats.CmpAscendingOrder(1.0, 0.0) < 0)
	assert.True(t, floats.CmpAscendingOrder(0.0, 0.0) == 0)

	assert.False(t, floats.CmpAscendingOrder(1.0, -0.5) > 0)
	assert.True(t, floats.CmpAscendingOrder(-0.5, 1.0) > 0)
	assert.False(t, floats.CmpAscendingOrder(-0.5, 1.0) < 0)
	assert.True(t, floats.CmpAscendingOrder(1.0, -0.5) < 0)
	assert.True(t, floats.CmpAscendingOrder(-0.5, -0.5) == 0)
}
