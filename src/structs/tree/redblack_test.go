package tree_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/structs/tree"
	"github.com/stretchr/testify/assert"
)

func TestRedBlackTree(t *testing.T) {
	tree := tree.NewRedBlack[float64, string](func(a, b float64) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		} else {
			return 0
		}
	})

	_, ok := tree.Left()
	assert.False(t, ok)

	_, ok = tree.Right()
	assert.False(t, ok)

	tree.Put(1.0, "1.0")
	if tree.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", tree.Size())
	}

	val, ok := tree.Get(1.0)
	if !ok {
		t.Errorf("Expected to find 1.0")
	} else if val != "1.0" {
		t.Errorf("Expected value to be 1.0, got %s", val)
	}

	left, ok := tree.Left()
	assert.True(t, ok)
	assert.Equal(t, "1.0", left)

	right, ok := tree.Right()
	assert.True(t, ok)
	assert.Equal(t, "1.0", right)

	tree.Remove(1.0)
	if tree.Size() != 0 {
		t.Errorf("Expected size to be 0, got %d", tree.Size())
	}

	val, ok = tree.Get(1.0)
	if ok {
		t.Errorf("Expected not to find 1.0")
	} else if val != "" {
		t.Errorf("Expected value to be empty, got %s", val)
	}

	tree.Remove(1.0)

	keys := tree.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected keys to be empty, got %v", keys)
	}

	values := tree.Values()
	if len(values) != 0 {
		t.Errorf("Expected values to be empty, got %v", values)
	}

	_, ok = tree.Left()
	assert.False(t, ok)

	_, ok = tree.Right()
	assert.False(t, ok)

	tree.Put(1.0, "1.0")
	tree.Put(2.0, "2.0")
	tree.Put(3.0, "3.0")

	keys = tree.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected keys to be [1.0, 2.0, 3.0], got %v", keys)
	} else if keys[0] != 1.0 {
		t.Errorf("Expected first key to be 1.0, got %v", keys[0])
	} else if keys[1] != 2.0 {
		t.Errorf("Expected second key to be 2.0, got %v", keys[1])
	} else if keys[2] != 3.0 {
		t.Errorf("Expected third key to be 3.0, got %v", keys[2])
	}

	values = tree.Values()
	if len(values) != 3 {
		t.Errorf("Expected values to be [1.0, 2.0, 3.0], got %v", values)
	} else if values[0] != "1.0" {
		t.Errorf("Expected first value to be 1.0, got %v", values[0])
	} else if values[1] != "2.0" {
		t.Errorf("Expected second value to be 2.0, got %v", values[1])
	} else if values[2] != "3.0" {
		t.Errorf("Expected third value to be 3.0, got %v", values[2])
	}

	left, ok = tree.Left()
	assert.True(t, ok)
	assert.Equal(t, "1.0", left)

	right, ok = tree.Right()
	assert.True(t, ok)
	assert.Equal(t, "3.0", right)

	tree.Put(1.5, "1.5")
	if tree.Size() != 4 {
		t.Errorf("Expected size to be 4, got %d", tree.Size())
	}

	val, ok = tree.Get(1.5)
	if !ok {
		t.Errorf("Expected to find 1.5")
	} else if val != "1.5" {
		t.Errorf("Expected value to be 1.5, got %s", val)
	}

	keys = tree.Keys()
	if len(keys) != 4 {
		t.Errorf("Expected keys to be [1.0, 1.5, 2.0, 3.0], got %v", keys)
	} else if keys[0] != 1.0 {
		t.Errorf("Expected first key to be 1.0, got %v", keys[0])
	} else if keys[1] != 1.5 {
		t.Errorf("Expected second key to be 1.5, got %v", keys[1])
	} else if keys[2] != 2.0 {
		t.Errorf("Expected third key to be 2.0, got %v", keys[2])
	} else if keys[3] != 3.0 {
		t.Errorf("Expected fourth key to be 3.0, got %v", keys[3])
	}

	values = tree.Values()
	if len(values) != 4 {
		t.Errorf("Expected values to be [1.0, 1.5, 2.0, 3.0], got %v", values)
	} else if values[0] != "1.0" {
		t.Errorf("Expected first value to be 1.0, got %v", values[0])
	} else if values[1] != "1.5" {
		t.Errorf("Expected second value to be 1.5, got %v", values[1])
	} else if values[2] != "2.0" {
		t.Errorf("Expected third value to be 2.0, got %v", values[2])
	} else if values[3] != "3.0" {
		t.Errorf("Expected fourth value to be 3.0, got %v", values[3])
	}

	left, ok = tree.Left()
	assert.True(t, ok)
	assert.Equal(t, "1.0", left)

	right, ok = tree.Right()
	assert.True(t, ok)
	assert.Equal(t, "3.0", right)

	val, ok = tree.Get(3.0)
	if !ok {
		t.Errorf("Expected to find 3.0")
	} else if val != "3.0" {
		t.Errorf("Expected value to be 3.0, got %s", val)
	}

	tree.Remove(4.0)
	if tree.Size() != 4 {
		t.Errorf("Expected size to be 4, got %d", tree.Size())
	}

	if tree.Empty() {
		t.Errorf("Expected tree not to be empty")
	}

	tree.Put(1.0, "2.0")
	if tree.Size() != 4 {
		t.Errorf("Expected size to be 4, got %d", tree.Size())
	}

	val, ok = tree.Get(1.0)
	if !ok {
		t.Errorf("Expected to find 1.0")
	} else if val != "2.0" {
		t.Errorf("Expected value to be 2.0, got %s", val)
	}

	tree.Clear()
	if tree.Size() != 0 {
		t.Errorf("Expected size to be 0, got %d", tree.Size())
	} else if !tree.Empty() {
		t.Errorf("Expected tree to be empty")
	}

	keys = tree.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected keys to be empty, got %v", keys)
	}

	values = tree.Values()
	if len(values) != 0 {
		t.Errorf("Expected values to be empty, got %v", values)
	}
}

func TestRedBlackTreeIterator(t *testing.T) {
	tree := tree.NewRedBlack[int, int](func(a, b int) int {
		return a - b
	})

	tree.Put(1, 1)
	tree.Put(2, 2)
	tree.Put(3, 3)

	it := tree.Iterator()
	assert.Panics(t, func() { it.Key() })
	assert.Panics(t, func() { it.Value() })

	// Must call Next before Key, Value, etc.
	assert.True(t, it.Next())
	assert.Equal(t, it.Key(), 1)
	assert.Equal(t, it.Value(), 1)

	assert.True(t, it.Next())
	assert.Equal(t, it.Key(), 2)
	assert.Equal(t, it.Value(), 2)

	assert.True(t, it.Next())
	assert.Equal(t, it.Key(), 3)
	assert.Equal(t, it.Value(), 3)

	// We can't call Key() or Value() if Next() returned false
	assert.False(t, it.Next())
	assert.Panics(t, func() { it.Key() })
	assert.Panics(t, func() { it.Value() })

	// Call it one more time to make sure it doesn't panic
	assert.False(t, it.Next())
	assert.True(t, it.Prev())

	// Reset the iterator
	it.Begin()
	assert.True(t, it.Next())
	assert.Equal(t, it.Key(), 1)
	assert.Equal(t, it.Value(), 1)

	// Set to the end of the iterator
	it.End()
	assert.False(t, it.Next())
	assert.Panics(t, func() { it.Key() })
	assert.Panics(t, func() { it.Value() })

	assert.True(t, it.Prev())
	assert.Equal(t, it.Key(), 3)
	assert.Equal(t, it.Value(), 3)

	assert.True(t, it.Prev())
	assert.Equal(t, it.Key(), 2)
	assert.Equal(t, it.Value(), 2)

	assert.True(t, it.Prev())
	assert.Equal(t, it.Key(), 1)
	assert.Equal(t, it.Value(), 1)

	assert.False(t, it.Prev())

	// Call it one more time to make sure it doesn't panic
	assert.False(t, it.Prev())
	assert.True(t, it.Next())

	// First and Last
	assert.True(t, it.First())
	assert.Equal(t, it.Key(), 1)
	assert.Equal(t, it.Value(), 1)
	assert.False(t, it.Prev())
	assert.True(t, it.Next())

	assert.True(t, it.Last())
	assert.Equal(t, it.Key(), 3)
	assert.Equal(t, it.Value(), 3)
	assert.False(t, it.Next())
	assert.True(t, it.Prev())
}
