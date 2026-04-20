package websocket

type SocketClient interface {
	Write(message []byte)
}
