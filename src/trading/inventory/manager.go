package inventory

import (
	"sync"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
)

type key struct {
	asset instr.Asset
	venue exchange.Exchange
}

type Update struct {
	Asset    instr.Asset
	Exchange exchange.Exchange
	Qty      float64
}

func (i *Update) key() key {
	return key{i.Asset, i.Exchange}
}

type Manager interface {
	Get(instr.Asset, exchange.Exchange) (float64, bool)
	Set(Update)
	Update(Update)
}

type SimpleManager struct {
	cache map[key]float64
}

func (i *SimpleManager) Get(asset instr.Asset, venue exchange.Exchange) (float64, bool) {
	qty, ok := i.cache[key{asset, venue}]
	return qty, ok
}

func (i *SimpleManager) Set(update Update) {
	i.cache[update.key()] = update.Qty
}

func (i *SimpleManager) Update(update Update) {
	i.cache[update.key()] += update.Qty
}

func NewSimpleManager() *SimpleManager {
	return &SimpleManager{
		cache: make(map[key]float64),
	}
}

type AtomicManager struct {
	cache *sync.Map
}

func (i *AtomicManager) Get(asset instr.Asset, venue exchange.Exchange) (float64, bool) {
    if val, ok := i.cache.Load(key{asset, venue}); ok {
        return val.(float64), true
    }
    return 0, false
}

func (i *AtomicManager) Set(update Update) {
	i.cache.Store(update.key(), update.Qty)
}

func (i *AtomicManager) Update(update Update) {
	for {
		if val, ok := i.cache.LoadOrStore(update.key(), update.Qty); !ok {
			return
		} else if i.cache.CompareAndSwap(update.key(), val, val.(float64)+update.Qty) {
			return
		}
	}
}

func NewAtomicManager() *AtomicManager {
	return &AtomicManager{
		cache: &sync.Map{},
	}
}
