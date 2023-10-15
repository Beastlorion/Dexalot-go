package orders

import (
	"strings"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/google/uuid"
)

func GenerateClientOrderID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

type Limit[T instr.Instrument] struct {
	Instrument       T
	ClientOrderID    string // Our ID
	ExchangeMetaData *exchange.MetaData
	Price            float64
	Qty              float64
	Side             side.OrderSide
	TimeInForce      tif.TimeInForce
	TransactTime     time.Time
}

type Market[T instr.Instrument] struct {
	Instrument       T
	ExchangeMetaData *exchange.MetaData
	ClientOrderID    string // Our ID
	Qty              float64
	Side             side.OrderSide
	TransactTime     time.Time
}
