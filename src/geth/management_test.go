package geth_test

import (
	"context"
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/geth"
	"github.com/ethereum/go-ethereum/core/types"
)

// Note - this test relies on the network
func TestManageHeaderConnection_AVAX_FUJI_C_CHAIN(t *testing.T) {
	avaxFujiWs := "wss://api.avax-test.network/ext/bc/C/ws"
	ctx, cancel := context.WithCancel(context.Background())
	header := make(chan *types.Header)
	manager := geth.NewWSConnectionManager(avaxFujiWs)

	start := time.Now()
	go func() {
		defer close(header)
		manager.ManageHeaderConnection(ctx, header)
	}()

	timeout := 30 * time.Second
	timer := time.NewTimer(timeout)
	select {
	case <-header:
		t.Logf("Received header after %f seconds", time.Since(start).Seconds())
		cancel()
	case <-timer.C:
		t.Errorf("Did not receive header after %f seconds", timeout.Seconds())
		cancel()
		t.Fail()
	}
}
