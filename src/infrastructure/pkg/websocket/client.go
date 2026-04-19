package websocket

import (
	"context"
	"everlasting/src/infrastructure/pkg/logger"
	"github.com/coder/websocket"
)

type WebSocketClient struct {
	client *websocket.Conn
	logger *logger.AppLogger
}

func NewWebSocketClient(client *websocket.Conn, logger *logger.AppLogger) *WebSocketClient {
	return &WebSocketClient{
		client: client,
		logger: logger,
	}
}

func (s *WebSocketClient) Write(ctx context.Context, message []byte) error {
	if err := s.client.Write(ctx, websocket.MessageText, message); err != nil {
		s.logger.Error(ctx, "websocket_client", err.Error())
		return err
	}
	return nil
}
