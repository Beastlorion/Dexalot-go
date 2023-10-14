package store_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/structs/store"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	store := store.New[int, string]()

	int1 := 1
	int2 := 2
	int3 := 3

	s1 := "s1"
	s2 := "s2"
	s3 := "s3"

	store.Store(int1, &s1)
	store.Store(int2, &s2)
	store.Store(int3, &s3)

	val, ok := store.Load(int1)
	assert.True(t, ok)
	assert.Equal(t, &s1, val)

	val, ok = store.Load(int2)
	assert.True(t, ok)
	assert.Equal(t, &s2, val)

	val, ok = store.Load(int3)
	assert.True(t, ok)
	assert.Equal(t, &s3, val)

	int4 := 4
	val, ok = store.Load(int4)
	assert.False(t, ok)
	assert.Nil(t, val)

	s1 = "new_s1"
	store.Store(int1, &s1)

	val, ok = store.Load(int1)
	assert.True(t, ok)
	assert.Equal(t, &s1, val)

	s4 := "new_s4"
	store.Store(int4, &s4)

	val, ok = store.Load(int4)
	assert.True(t, ok)
	assert.Equal(t, &s4, val)

	N := 100_000

	wg := &sync.WaitGroup{}
	wg.Add(6)

	ctx1, cancel1 := context.WithCancel(context.Background())
	go func(cancel context.CancelFunc) {
		defer wg.Done()
		for i := 0; i < N; i++ {
			s := strconv.Itoa(i)
			store.Store(int1, &s)
		}
		time.Sleep(time.Millisecond)
		cancel()
	}(cancel1)

	ctx2, cancel2 := context.WithCancel(context.Background())
	go func(cancel context.CancelFunc) {
		defer wg.Done()
		for i := 0; i < N; i++ {
			s := strconv.Itoa(i)
			store.Store(int2, &s)
		}
		time.Sleep(time.Millisecond)
		cancel()
	}(cancel2)

	ctx3, cancel3 := context.WithCancel(context.Background())
	go func(cancel context.CancelFunc) {
		defer wg.Done()
		for i := 0; i < N; i++ {
			s := strconv.Itoa(i)
			store.Store(int3, &s)
		}
		time.Sleep(time.Millisecond)
		cancel()
	}(cancel3)

	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				val, ok := store.Load(int1)
				assert.True(t, ok)
				assert.NotNil(t, val)
			}
		}
	}(ctx1)

	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				val, ok := store.Load(int2)
				assert.True(t, ok)
				assert.NotNil(t, val)
			}
		}
	}(ctx2)

	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				val, ok := store.Load(int3)
				assert.True(t, ok)
				assert.NotNil(t, val)
			}
		}
	}(ctx3)

	wg.Wait()

	val, ok = store.Load(int1)
	assert.True(t, ok)
	assert.NotNil(t, val)
	parsedInt, _ := strconv.ParseInt(*val, 10, 64)
	assert.Equal(t, int64(N-1), parsedInt)

	val, ok = store.Load(int2)
	assert.True(t, ok)
	assert.NotNil(t, val)
	parsedInt, _ = strconv.ParseInt(*val, 10, 64)
	assert.Equal(t, int64(N-1), parsedInt)

	val, ok = store.Load(int3)
	assert.True(t, ok)
	assert.NotNil(t, val)
	parsedInt, _ = strconv.ParseInt(*val, 10, 64)
	assert.Equal(t, int64(N-1), parsedInt)
}
