package ticker

import (
	"math"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/source"
)

type Ticker[T instr.Instrument] struct {
	Instrument     T                 `json:"instrument"`
	Source         source.Source `json:"exchange"`
	Bid            float64           `json:"bid"`
	Offer          float64           `json:"offer"`
	Mid            float64           `json:"mid"`
	Last           float64           `json:"last"`
	SequenceNumber int64             `json:"sequenceNumber"`
	TransactTime   time.Time         `json:"transactTime"`
}

func EmptyTicker[T instr.Instrument]() Ticker[T] {
	return Ticker[T]{
		Bid:            math.NaN(),
		Offer:          math.NaN(),
		Mid:            math.NaN(),
		Last:           math.NaN(),
		SequenceNumber: 0,
		TransactTime:   time.Unix(0, 0),
	}
}
