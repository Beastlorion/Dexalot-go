package tif

type TimeInForce int

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
