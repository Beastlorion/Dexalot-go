package store

import "sync"

// Note: This will always grow by design. Having to delete data is a sign that we are doing something wrong.
type Store[T comparable, U any] struct {
	cache *sync.Map
}

func (m *Store[T, U]) Load(instrument T) (*U, bool) {
	if val, ok := m.cache.Load(instrument); ok {
		return val.(*U), true
	}
	return nil, false
}

func (m *Store[T, U]) Store(instrument T, data *U) {
	m.cache.Store(instrument, data)
}

func New[T comparable, U any]() *Store[T, U] {
	return &Store[T, U]{
		cache: &sync.Map{},
	}
}
