package report

import (
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/exchange"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordstatus"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordtypes"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
)

// Note: Optionals are indicated as pointers
type ExecutionReport[T instr.Instrument] struct {
	Instrument    T                 // Instrument of the order.
	Exchange      exchange.Exchange          // Exchange of the order.
	ClientOrderID string            // Unique ID for the order assigned by the client (us).
	OrderID       string            // Unique ID for the order assigned by the exchange.
	Price         float64           // Price of the order.
	Qty           float64           // Qty of the order.
	Side          side.OrderSide // Side of the order.
	OrderStatus   ordstatus.OrderStatus       // OrderStatus of the order.
	OrderType     ordtypes.OrderType         // Type of the order.
	TimeInForce   tif.TimeInForce       // Time in force of the order.
	TransactTime  *time.Time        // Transact time of the order
	Account       *string           // Account of the order (optional).
	TotalFee      *float64          // Total fee of the order (optional).
	FilledQty     *float64          // Qty of the order that has been filled.
	LeavesQty     *float64          // Leaves qty of the order (optional).
	CumQty        *float64          // Cumulative qty of the order (optional).
	AvgPrice      *float64          // Average price of the order (optional).
	LastPrice     *float64          // Last price of the order (optional).
	Version       *uint8            // Exchange version of the order (optional).
	RejectReason  *string           // Reason for the order being rejected (optional).
	ErrorReason   *string           // Reason for the order being errored (optional).
}

type Builder[T instr.Instrument] struct {
	report *ExecutionReport[T]
}

func NewBuilder[T instr.Instrument](
	instrument T,
	exchange exchange.Exchange,
	clientOrderID string,
	orderID string,
	price float64,
	qty float64,
	side side.OrderSide,
	orderStatus ordstatus.OrderStatus,
	orderType ordtypes.OrderType,
	timeInForce tif.TimeInForce) *Builder[T] {
	return &Builder[T]{
		report: &ExecutionReport[T]{
			Instrument:    instrument,
			Exchange:      exchange,
			ClientOrderID: clientOrderID,
			OrderID:       orderID,
			Price:         price,
			Qty:           qty,
			Side:          side,
			OrderStatus:   orderStatus,
			OrderType:     orderType,
			TimeInForce:   timeInForce,
		},
	}
}

func (b *Builder[T]) Build() *ExecutionReport[T] {
	return b.report
}

func (b *Builder[T]) SetTransactTime(t time.Time) *Builder[T] {
	b.report.TransactTime = &t
	return b
}

func (b *Builder[T]) SetAccount(account string) *Builder[T] {
	b.report.Account = &account
	return b
}

func (b *Builder[T]) SetTotalFee(totalFee float64) *Builder[T] {
	b.report.TotalFee = &totalFee
	return b
}

func (b *Builder[T]) SetFilledQty(filledQty float64) *Builder[T] {
	b.report.FilledQty = &filledQty
	return b
}

func (b *Builder[T]) SetLeavesQty(leavesQty float64) *Builder[T] {
	b.report.LeavesQty = &leavesQty
	return b
}

func (b *Builder[T]) SetCumQty(cumQty float64) *Builder[T] {
	b.report.CumQty = &cumQty
	return b
}

func (b *Builder[T]) SetAvgPrice(avgPrice float64) *Builder[T] {
	b.report.AvgPrice = &avgPrice
	return b
}

func (b *Builder[T]) SetLastPrice(lastPrice float64) *Builder[T] {
	b.report.LastPrice = &lastPrice
	return b
}

func (b *Builder[T]) SetVersion(version uint8) *Builder[T] {
	b.report.Version = &version
	return b
}

func (b *Builder[T]) SetRejectReason(rejectReason string) *Builder[T] {
	b.report.RejectReason = &rejectReason
	return b
}

func (b *Builder[T]) SetErrorReason(errorReason string) *Builder[T] {
	b.report.ErrorReason = &errorReason
	return b
}
