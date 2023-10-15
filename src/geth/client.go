package geth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewWebsocketClient(ctx context.Context, url string) (*ethclient.Client, error) {
	c, err := rpc.DialWebsocket(ctx, url, "")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %s", url, err.Error())
	}
	return ethclient.NewClient(c), nil
}

func NewRPCClient(url string) (*ethclient.Client, error) {
	c, err := rpc.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %s", url, err.Error())
	}
	return ethclient.NewClient(c), nil
}
