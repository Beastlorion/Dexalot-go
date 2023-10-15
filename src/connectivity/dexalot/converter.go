package dexalot

import (
	"fmt"
	"strings"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/instr"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/marketdata/types/side"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordstatus"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/ordtypes"
	"github.com/Abso1ut3Zer0/Dexalot-go/src/trading/constants/tif"
	"github.com/ethereum/go-ethereum/common"
)

type Method string
type SocketDataType string

const (
	Subscribe                   = Method("subscribe")
	ChartSubscribe              = Method("chartsubscribe")
	MarketDataSnapshotSubscribe = Method("marketsnapshotsubscribe")
)

const (
	OrderBooks     = SocketDataType("orderBooks")
	LastTrades     = SocketDataType("lastTrades")
	MarketSnapshot = SocketDataType("marketSnapShot")
	ChartSnapshot  = SocketDataType("chartSnapShot")
	Prices         = SocketDataType("prices")
	AuctionData    = SocketDataType("auctionData")
)

type OrderStatus int

const (
	New = OrderStatus(iota)
	Rejected
	PartiallyFilled
	Filled
	Canceled
	Expired
	Killed
	CancelReject
)

type OrderType int

const (
	Market = iota
	Limit
	Stop
	StopLimit
)

type TimeInForce int

const (
	GTC = iota
	FOK
	IOC
	PO
)

func ConvertSideToAPI(s side.OrderSide) int64 {
	if s == side.BUY {
		return 0
	}
	return 1
}

func ConvertToSide(s uint8) (side.OrderSide, error) {
	if s == 0 {
		return side.BUY, nil
	} else if s == 1 {
		return side.SELL, nil
	}
	return side.NilSide, fmt.Errorf("unknown side: %d", s)
}

func ConvertFromAPIAsset(asset string) instr.Asset {
	return instr.Asset(strings.ToUpper(asset))
}

func ConvertToAPIAsset(asset instr.Asset) (string, error) {
	switch asset {
	case instr.USDT:
		return "USDt", nil
	default:
		return string(asset), nil
	}
}

func ConvertSpotToAPIInstrument(instrument instr.Spot) (string, error) {
	base, err := ConvertToAPIAsset(instrument.Base)
	if err != nil {
		return "", err
	}
	term, err := ConvertToAPIAsset(instrument.Term)
	if err != nil {
		return "", err
	}
	return base + "/" + term, nil
}

func ConvertToSpot(name string) (instr.Spot, error) {
	return instr.ToSpotSlashSeparated(name)
}

func ConvertFromAPIOrderStatus(status OrderStatus) (ordstatus.OrderStatus, error) {
	switch status {
	case New:
		return ordstatus.New, nil
	case Rejected:
		return ordstatus.Rejected, nil
	case PartiallyFilled:
		return ordstatus.PartiallyFilled, nil
	case Filled:
		return ordstatus.Filled, nil
	case Canceled:
		return ordstatus.Canceled, nil
	case Expired:
		return ordstatus.Expired, nil
	case CancelReject:
		return ordstatus.CancelReject, nil
	}
	return ordstatus.NilStatus, fmt.Errorf("unknown order status: %d", status)
}

func ConvertToAPIOrderType(orderType ordtypes.OrderType) (OrderType, error) {
	switch orderType {
	case ordtypes.Market:
		return Market, nil
	case ordtypes.Limit:
		return Limit, nil
	case ordtypes.Stop:
		return Stop, nil
	case ordtypes.StopLimit:
		return StopLimit, nil
	}
	return -1, fmt.Errorf("unknown order type: %d", orderType)
}

func ConvertFromAPIOrderType(orderType OrderType) (ordtypes.OrderType, error) {
	switch orderType {
	case Market:
		return ordtypes.Market, nil
	case Limit:
		return ordtypes.Limit, nil
	case Stop:
		return ordtypes.Stop, nil
	case StopLimit:
		return ordtypes.StopLimit, nil
	}
	return ordtypes.NilOrderType, fmt.Errorf("unknown order type: %d", orderType)
}

func ConvertToAPITimeInForce(timeInForce tif.TimeInForce) (TimeInForce, error) {
	switch timeInForce {
	case tif.GTC:
		return GTC, nil
	case tif.FOK:
		return FOK, nil
	case tif.IOC:
		return IOC, nil
	case tif.PO:
		return PO, nil
	}
	return -1, fmt.Errorf("unknown time in force: %d", timeInForce)
}

func ConvertFromAPITimeInForce(timeInForce TimeInForce) (tif.TimeInForce, error) {
	switch timeInForce {
	case GTC:
		return tif.GTC, nil
	case FOK:
		return tif.FOK, nil
	case IOC:
		return tif.IOC, nil
	case PO:
		return tif.PO, nil
	}
	return tif.NilTimeInForce, fmt.Errorf("unknown time in force: %d", timeInForce)
}

func CreateSpotHash(instrument instr.Spot) (common.Hash, error) {
	instrString, err := ConvertSpotToAPIInstrument(instrument)
	if err != nil {
		return common.Hash{}, err
	}
	var instrByte32 [32]byte
	copy(instrByte32[:], instrString)
	return common.BytesToHash(instrByte32[:]), nil
}
