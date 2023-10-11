package trtypes

import (
	"strings"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/mdtypes"
	"github.com/google/uuid"
)

func GenerateClientOrderID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

type LimitOrder[T instr.Instrument] struct {
	Instrument       T
	ClientOrderID    string // Our ID
	ExchangeMetaData *ExchangeMetaData
	Price            float64
	Qty              float64
	Side             mdtypes.OrderSide
	TimeInForce      TimeInForce
	TransactTime     time.Time
}

type MarketOrder[T instr.Instrument] struct {
	Instrument       T
	ExchangeMetaData *ExchangeMetaData
	ClientOrderID    string // Our ID
	Qty              float64
	Side             mdtypes.OrderSide
	TransactTime     time.Time
}

