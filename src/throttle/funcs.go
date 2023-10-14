package throttle

import (
	"context"
	"math/rand"
	"time"
)

const (
	defaultMinPeriod = 25 * time.Millisecond
	defaultMaxPeriod = 500 * time.Millisecond
)

func DefaultRandom() {
	RandomRange(defaultMinPeriod, defaultMaxPeriod)
}

func RandomRange(min, max time.Duration) {
	time.Sleep(time.Duration(rand.Int63n(int64(max-min))) + min)
}

func DefaultRandomCall(ctx context.Context, call func(context.Context) error) error {
	return RandomCallRange(ctx, defaultMinPeriod, defaultMaxPeriod, call)
}

func RandomCallRange(ctx context.Context, min, max time.Duration, call func(context.Context) error) error {
	RandomRange(min, max)
	return call(ctx)
}
