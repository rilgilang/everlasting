package event

import (
	"context"
	"encoding/json"
	"everlasting/src/domain/event"
	"everlasting/src/infrastructure/pkg"
	"everlasting/src/infrastructure/pkg/logger"
	"fmt"
	"github.com/sarulabs/di"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
	"time"
)

type WishingWallMessage struct {
	container di.Container
	config    *pkg.Config
}

func (u *WishingWallMessage) Handle(ctx context.Context, body []byte) (err error) {
	var (
		logger = u.container.Get("logger.app").(*logger.AppLogger)
	)

	messageRequest := new(event.WishingWallMessage)
	err = json.Unmarshal(body, messageRequest)
	if err != nil {
		logger.Error(ctx, "wishing_wall_message", err.Error())
		return err
	}

	fmt.Println("Processing message --> ", messageRequest)

	// Connect to the SAME namespace as Postman (root namespace "/")
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["event"] = messageRequest.EventID

	// Use the correct URI - root namespace
	uri := "http://127.0.0.1:8000/"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		logger.Error(ctx, "wishing_wall_message_connect", err.Error())
		return err
	}
	//defer client.Close()

	// Wait for connection to be established
	connected := make(chan bool)

	client.On("connect", func() {
		fmt.Println("Consumer connected to socket server")
		connected <- true
	})

	client.On("error", func(err error) {
		logger.Error(ctx, "wishing_wall_message_socket_error", err.Error())
	})

	// Wait for connection or timeout
	select {
	case <-connected:
		// Connection established, send message
		fmt.Println("Sending message to room:", messageRequest.EventID)
		client.Emit("msg", messageRequest.Message)

		// Give time for the message to be sent
		// You might want to add a small delay
		// time.Sleep(100 * time.Millisecond)

	case <-time.After(5 * time.Second):
		logger.Error(ctx, "wishing_wall_message_timeout", "Connection timeout")
		return fmt.Errorf("connection timeout")
	}

	fmt.Println("Message sent successfully")
	return nil
}

func NewWishingWallMessage(container di.Container, config *pkg.Config) *WishingWallMessage {
	return &WishingWallMessage{container: container, config: config}
}
