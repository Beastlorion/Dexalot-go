package trtypes

type OrderStatus int
type OrderType int
type TimeInForce int

const (
	NilStatus = OrderStatus(iota)
	New
	PartiallyFilled
	Filled
	DoneForDay
	Canceled
	Replaced
	PendingCancel
	Stopped
	Rejected
	Suspended
	PendingNew
	Calculated
	Expired
	AcceptedForBidding
	PendingReplace
	CancelReject
	ErrorStatus // Internal error
)

const (
	NilOrderType = OrderType(iota)
	Market
	Limit
	Stop
	StopLimit
)

const (
	NilTimeInForce = TimeInForce(iota)
	DAY            // Good Till Day
	GTC            // Good Till CancelBatchLimitOrders
	ATO            // At the Opening
	IOC            // Immediate or CancelBatchLimitOrders
	FOK            // Fill or Kill
	GTX            // Good Till Crossing
	GTD            // Good Till Date
	PO             // Post Only
)
