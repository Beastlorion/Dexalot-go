package geth_test

 import (
	"context"
	"testing"
	"time"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/geth"
	"github.com/ethereum/go-ethereum/core/types"
)

type testHandler struct {
	cancel context.CancelFunc
	t      *testing.T
}

func (h *testHandler) Handle(logs []types.Log) {
	for _, log := range logs {
		h.t.Logf("Received log from block: %d", log.BlockNumber)
	}
	h.cancel()
}

// Note: this test relies on the network
func TestReceiveHeaders_AVAX_FUJI_C_CHAIN(t *testing.T) {
	avaxFujiWs := "wss://api.avax-test.network/ext/bc/C/ws"
	avaxFujiRPC := "https://api.avax-test.network/ext/bc/C/rpc"
	ctx, cancel := context.WithCancel(context.Background())
	header := make(chan *types.Header)

	start := time.Now()
	rpcClient, err := geth.NewRPCClient(avaxFujiRPC)
	if err != nil {
		t.Fatalf("failed to connect to %s: %s", avaxFujiRPC, err.Error())
	}
	handlerCtx, handlerCancel := context.WithCancel(context.Background())
	handler := &testHandler{cancel: handlerCancel, t: t}
	reciever := geth.NewHeaderReceiver(rpcClient, 30 * time.Second)
	go reciever.Receive(ctx, handler, header)
	
	manager := geth.NewWSConnectionManager(avaxFujiWs)
	go func() {
		defer close(header)
		manager.ManageHeaderConnection(ctx, header)
	}()

	timeout := 30 * time.Second
	timer := time.NewTimer(timeout)
	select {
	case <-handlerCtx.Done():
		t.Logf("Received header after %f seconds", time.Since(start).Seconds())
		cancel()
	case <-timer.C:
		t.Errorf("Did not receive header after %f seconds", timeout.Seconds())
		cancel()
		handlerCancel()
		t.Fail()
	}
}
