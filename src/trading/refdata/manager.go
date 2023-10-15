package refdata

import (
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
)

type managerKey[T instr.Instrument] struct {
	Exchange   exchange.Exchange
	Instrument T
}

type Manager[T instr.Instrument] struct {
	refData map[managerKey[T]]*Composite
}

func (r *Manager[T]) Register(exchange exchange.Exchange, instrument T, refData *Composite) {
	r.refData[managerKey[T]{Exchange: exchange, Instrument: instrument}] = refData
}

func (r *Manager[T]) Get(exchange exchange.Exchange, instrument T) (*Composite, error) {
	if data, ok := r.refData[managerKey[T]{Exchange: exchange, Instrument: instrument}]; ok {
		return data, nil
	}
	return nil, ReferenceDataNotFoundError
}

func NewManager[T instr.Instrument]() *Manager[T] {
	return &Manager[T]{
		refData: make(map[managerKey[T]]*Composite),
	}
}
