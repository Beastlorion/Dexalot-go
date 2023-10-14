package ordstatus

type OrderStatus int

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
