package trtypes

import (
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/mdtypes"
)

// Note: Optionals are indicated as pointers
type ExecutionReport[T instr.Instrument] struct {
	Instrument    T                    // Instrument of the order.
	Exchange      Exchange   // Exchange of the order.
	ClientOrderID string               // Unique ID for the order assigned by the client (us).
	OrderID       string               // Unique ID for the order assigned by the exchange.
	Price         float64              // Price of the order.
	Qty           float64              // Qty of the order.
	Side          mdtypes.OrderSide // Side of the order.
	OrderStatus   OrderStatus          // OrderStatus of the order.
	OrderType     OrderType            // Type of the order.
	TimeInForce   TimeInForce          // Time in force of the order.
	TransactTime  *time.Time            // Transact time of the order
	Account       *string               // Account of the order (optional).
	TotalFee      *float64              // Total fee of the order (optional).
	FilledQty     *float64              // Qty of the order that has been filled.
	LeavesQty     *float64              // Leaves qty of the order (optional).
	CumQty        *float64              // Cumulative qty of the order (optional).
	AvgPrice      *float64              // Average price of the order (optional).
	LastPrice     *float64              // Last price of the order (optional).
	Version       *uint8                // Exchange version of the order (optional).
	RejectReason  *string               // Reason for the order being rejected (optional).
	ErrorReason   *string               // Reason for the order being errored (optional).
}

