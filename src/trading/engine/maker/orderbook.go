package maker

import (
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/structs/tree"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordstatus"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordtypes"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/orders"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/report"
)

type Layer struct {
	ClientOrderID string
	Price         float64
}

type OrderBook[T instr.Instrument] struct {
	instrument T
	exchange   exchange.Exchange
	Bids       *tree.RedBlack[Layer, *orders.Limit[T]]
	Offers     *tree.RedBlack[Layer, *orders.Limit[T]]
}

func (book *OrderBook[T]) ApplyExecutionReport(rep *report.ExecutionReport[T]) {
	if rep.Instrument != book.instrument || rep.Exchange != book.exchange || rep.OrderType != ordtypes.Limit {
		return
	}

	if rep.TimeInForce != tif.GTC && rep.TimeInForce != tif.PO {
		return
	}

	switch rep.OrderStatus {
	case ordstatus.New:
		order := &orders.Limit[T]{
			Instrument:       rep.Instrument,
			ExchangeMetaData: &exchange.MetaData{Exchange: rep.Exchange, OrderID: &rep.OrderID},
			ClientOrderID:    rep.ClientOrderID,
			Price:            rep.Price,
			Qty:              rep.Qty,
			Side:             rep.Side,
			TimeInForce:      rep.TimeInForce,
			TransactTime:     time.Now(),
		}

		switch rep.Side {
		case side.BUY:
			book.Bids.Put(Layer{rep.ClientOrderID, rep.Price}, order)
		case side.SELL:
			book.Offers.Put(Layer{rep.ClientOrderID, rep.Price}, order)
		}
	case ordstatus.Filled:
		fallthrough
	case ordstatus.DoneForDay:
		fallthrough
	case ordstatus.Replaced:
		fallthrough
	case ordstatus.Expired:
		fallthrough
	case ordstatus.Canceled:
		switch rep.Side {
		case side.BUY:
			book.Bids.Remove(Layer{rep.ClientOrderID, rep.Price})
		case side.SELL:
			book.Offers.Remove(Layer{rep.ClientOrderID, rep.Price})
		}
	}
}

func NewOrderBook[T instr.Instrument](instrument T, exchange exchange.Exchange) *OrderBook[T] {
	return &OrderBook[T]{
		instrument: instrument,
		exchange:   exchange,
		Bids: tree.NewRedBlack[Layer, *orders.Limit[T]](func(a, b Layer) int {
			if a.Price > b.Price {
				return -1
			} else if a.Price < b.Price {
				return 1
			}

			if a.ClientOrderID < b.ClientOrderID {
				return -1
			} else if a.ClientOrderID > b.ClientOrderID {
				return 1
			}
			return 0
		}),
		Offers: tree.NewRedBlack[Layer, *orders.Limit[T]](func(a, b Layer) int {
			if a.Price < b.Price {
				return -1
			} else if a.Price > b.Price {
				return 1
			}

			if a.ClientOrderID < b.ClientOrderID {
				return -1
			} else if a.ClientOrderID > b.ClientOrderID {
				return 1
			}

			return 0
		}),
	}
}
