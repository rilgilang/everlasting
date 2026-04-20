package websocket

import (
	"context"
	"everlasting/src/infrastructure/pkg/logger"
	"github.com/coder/websocket"
)

type WebSocketClient struct {
	client *websocket.Conn
	send   chan []byte
	logger *logger.AppLogger
}

func NewWebSocketClient(client *websocket.Conn, logger *logger.AppLogger) *WebSocketClient {
	conn := &WebSocketClient{
		client: client,
		send:   make(chan []byte, 100),
		logger: logger,
	}

	go conn.writeLoop(context.Background())

	return conn
}

func (c *WebSocketClient) Write(message []byte) {
	select {
	case c.send <- message:
	default:
		c.logger.Error(context.Background(), "websocket_client", "send buffer full")
	}
}

func (c *WebSocketClient) writeLoop(ctx context.Context) {
	for {
		select {
		case msg := <-c.send:
			err := c.client.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				c.logger.Error(ctx, "websocket_write", err.Error())
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
