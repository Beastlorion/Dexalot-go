package ws

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const reconnectTime = 250 * time.Millisecond


type Client struct {
	c         *websocket.Conn
	onMessage func([]byte)
}

// Note: this is required to fulfill the websocket protocol.
func (client *Client) SetHandlers(pingWriteTimeout, pingReadTimeout time.Duration) {
	client.c.SetPingHandler(func(data string) error {
		return client.c.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(pingWriteTimeout))
	})

	client.c.SetPongHandler(func(data string) error {
		return client.c.SetReadDeadline(time.Now().Add(pingReadTimeout))
	})
}

func (client *Client) Connect(url url.URL, header http.Header) error {
	c, res, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("failed to connect, status code: %d", res.StatusCode)
	}
	client.c = c
	return nil
}

func (client *Client) Subscribe(handshake map[string]any) error {
	err := client.c.WriteJSON(handshake)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Read(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, b, err := client.c.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				} else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Warn("Unexpected error reading from websocket: %v", err)
					return
				} else {
					slog.Warn("Error reading from websocket: %v", err)
					return
				}
			}
			client.onMessage(b)
		}
	}
}

func (client *Client) Ping(ctx context.Context, cancel context.CancelFunc, pingInterval, pingWriteTimeout time.Duration) {
	defer cancel()
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			err := client.c.WriteControl(websocket.PingMessage, nil, time.Now().Add(pingWriteTimeout))
			if err != nil {
				slog.Warn("Did not receive pong from server")
				return
			}
		}
	}
}

func (client *Client) Close() error {
	return client.c.Close()
}

func New(onMessage func([]byte)) *Client {
	return &Client{
		onMessage: onMessage,
	}
}
