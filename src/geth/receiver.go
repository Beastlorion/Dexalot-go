package geth

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type HeaderProcessor interface {
	ProcessHead(ctx context.Context, header *types.Header) error
}

// this is an area of improvement since we should have knowledge about what logs are relevant
type EventHandler interface {
	Handle([]types.Log)
}

// this is an area of improvement, as we could handle multiple filters
type HeaderReceiver struct {
	c       *ethclient.Client
	timeout time.Duration
	filter  ethereum.FilterQuery
}

func NewHeaderReceiver(c *ethclient.Client, timeout time.Duration) *HeaderReceiver {
	return &HeaderReceiver{
		c:       c,
		timeout: timeout,
		filter:  ethereum.FilterQuery{}, // we will just get all logs for a block range
	}
}

func (h *HeaderReceiver) Receive(ctx context.Context, handler EventHandler, header <-chan *types.Header) {
	for {
		select {
		case <-ctx.Done():
			return
		case head, ok := <-header:
			if !ok {
				return // channel closed
			}

			if h.filter.FromBlock == nil {
				h.filter.FromBlock = head.Number
			}
			h.filter.ToBlock = head.Number

			tCtx, cancel := context.WithTimeout(ctx, h.timeout)
			logs, err := h.c.FilterLogs(tCtx, h.filter)
			cancel() // release context resources
			if err != nil {
				slog.Warn("Could not fetch logs", "error", err.Error())
				continue
			}

			handler.Handle(logs)
			h.filter.FromBlock.Add(head.Number, big.NewInt(1))
		}
	}
}
