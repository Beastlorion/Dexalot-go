package tree

import "github.com/emirpasic/gods/trees/redblacktree"

type NodeIterator[K comparable, V any] struct {
	iter redblacktree.Iterator
}

func (i *NodeIterator[K, V]) Begin() {
	i.iter.Begin()
}

// Note: you can immediately call Key() and Value() after calling First()
func (i *NodeIterator[K, V]) First() bool {
	return i.iter.First()
}

func (i *NodeIterator[K, V]) End() {
	i.iter.End()
}

// Note: you can immediately call Key() and Value() after calling Last()
func (i *NodeIterator[K, V]) Last() bool {
	return i.iter.Last()
}

// Note: this must be called first if we are at the beginning, otherwise it will panic
func (i *NodeIterator[K, V]) Next() bool {
	return i.iter.Next()
}

// Note: this must be called first if we are at the end, otherwise it will panic
func (i *NodeIterator[K, V]) Prev() bool {
	return i.iter.Prev()
}

func (i *NodeIterator[K, V]) Key() K {
	if key, ok := i.iter.Key().(K); ok {
		return key
	}
	var zero K
	return zero
}

func (i *NodeIterator[K, V]) Value() V {
	if val, ok := i.iter.Value().(V); ok {
		return val
	}
	var zero V
	return zero
}

type RedBlack[K comparable, V any] struct {
	tree *redblacktree.Tree
	iter *NodeIterator[K, V]
}

func (t *RedBlack[K, V]) Put(key K, value V) {
	t.tree.Put(key, value)
}

func (t *RedBlack[K, V]) Get(key K) (V, bool) {
	if v, ok := t.tree.Get(key); ok {
		return v.(V), true
	}
	var zero V
	return zero, false
}

func (t *RedBlack[K, V]) Remove(key K) {
	t.tree.Remove(key)
}

func (t *RedBlack[K, V]) Empty() bool {
	return t.tree.Empty()
}

func (t *RedBlack[K, V]) Size() int {
	return t.tree.Size()
}

func (t *RedBlack[K, V]) Keys() []K {
	keys := make([]K, 0, t.tree.Size())
	for _, key := range t.tree.Keys() {
		keys = append(keys, key.(K))
	}
	return keys
}

func (t *RedBlack[K, V]) Values() []V {
	values := make([]V, 0, t.tree.Size())
	for _, value := range t.tree.Values() {
		values = append(values, value.(V))
	}
	return values
}

func (t *RedBlack[K, V]) Clear() {
	t.tree.Clear()
}

func (t *RedBlack[K, V]) Left() (V, bool) {
	if t.Empty() {
		var zero V
		return zero, false
	}
	return t.tree.Left().Value.(V), true
}

func (t *RedBlack[K, V]) Right() (V, bool) {
	if t.Empty() {
		var zero V
		return zero, false
	}
	return t.tree.Right().Value.(V), true
}

func (t *RedBlack[K, V]) Iterator() NodeIterator[K, V] {
	return NodeIterator[K, V]{iter: t.tree.Iterator()}
}

// TODO - rely on custom red black tree implementation
func NewRedBlack[K comparable, V any](comparator func(a, b K) int) *RedBlack[K, V] {
	return &RedBlack[K, V]{
		tree: redblacktree.NewWith(func(a, b interface{}) int {
			return comparator(a.(K), b.(K))
		}),
	}
}
