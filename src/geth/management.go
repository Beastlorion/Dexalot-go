package geth

import (
	"context"
	"log/slog"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/throttle"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	maxThrottlePeriod       = 250 * time.Millisecond
	connectionRetryInterval = 5 * time.Second
	heartbeatInterval       = 30 * time.Second
	connectionTimeout       = 2 * time.Second
)

type WSConnectionManager struct {
	wsURL    string
	c        *ethclient.Client
	isActive bool
}

func NewWSConnectionManager(wsURL string) *WSConnectionManager {
	return &WSConnectionManager{
		wsURL: wsURL,
	}
}

func (m *WSConnectionManager) watchHeaderSubscription(ctx context.Context, sub ethereum.Subscription, header chan<- *types.Header) {
	heartBeat := time.NewTicker(heartbeatInterval)
	for {
		select {
		case <-ctx.Done():
			heartBeat.Stop()
			sub.Unsubscribe()
			m.c.Close()
			m.isActive = false
			return
		case <-heartBeat.C:
			slog.Info("Node Heartbeat", "pending", len(header), "capacity", cap(header))
		case err := <-sub.Err():
			slog.Warn("Connection error", "error", err.Error())
			heartBeat.Stop()
			sub.Unsubscribe()
			slog.Info("Reconnecting to node")
			time.Sleep(connectionRetryInterval)
		}
	}
}

func (m *WSConnectionManager) ManageHeaderConnection(ctx context.Context, header chan<- *types.Header) {
	m.isActive = true
	i := 0
	for m.isActive {
		throttle.DefaultRandom() // prevent connection stampede
		connCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
		conn, err := NewWebsocketClient(connCtx, m.wsURL)
		if err != nil {
			slog.Warn("Could not connect to node", "error", err.Error())
			slog.Info("Reattemtping connection to node")
			time.Sleep(connectionRetryInterval << uint(i))
			i++
			continue
		}
		m.c = conn

		subCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
		sub, err := conn.SubscribeNewHead(subCtx, header)
		cancel()
		if err != nil {
			slog.Warn("Could not subscribe to node", "error", err.Error())
			slog.Info("Reattemtping connection to node")
			conn.Close()
			time.Sleep(connectionRetryInterval << uint(i))
			i++
			continue
		}

		i = 0
		m.watchHeaderSubscription(ctx, sub, header)
	}
}
