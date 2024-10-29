package websocketproxy

import (
	"context"
	"net"
	"time"

	"nhooyr.io/websocket"
)

type WebSocketProxyClient struct {
	timeout time.Duration
}

func NewWebSocketProxyClient(timeout time.Duration) *WebSocketProxyClient {
	return &WebSocketProxyClient{timeout}
}

func (p *WebSocketProxyClient) Dialer(ctx context.Context, url string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return websocket.NetConn(context.Background(), conn, websocket.MessageBinary), nil
}
