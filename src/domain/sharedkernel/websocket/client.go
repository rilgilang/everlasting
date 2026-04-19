package websocket

import (
	"context"
)

type SocketClient interface {
	Write(ctx context.Context, message []byte) (err error)
}
