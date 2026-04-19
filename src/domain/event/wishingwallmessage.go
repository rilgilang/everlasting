package event

import (
	"context"
	"encoding/json"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/messagebroker"
	"everlasting/src/domain/sharedkernel/photo"
	minio "everlasting/src/domain/sharedkernel/storage"
	"everlasting/src/domain/sharedkernel/websocket"
	"fmt"
	"time"
)

const WishingWallMessageTask messagebroker.TaskName = "wishing_wall_message"

type (
	WishingWallInput struct {
		Name    string      `json:"name"`
		Message string      `json:"message"`
		EventID string      `json:"event_id"`
		Photo   photo.Photo `json:"photos"`
	}

	WishingWallMessage struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Message   string    `json:"message"`
		Photo     string    `json:"photo"`
		EventID   string    `json:"event_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (ww *WishingWallInput) SaveMessage(ctx context.Context, repository WishingWallMessageRepository, storage minio.StorageRepository, socketClient websocket.SocketClient) (*WishingWallMessage, error) {

	// Store the photo into Storage
	dir := fmt.Sprintf(`%s/%v.%s`, ww.EventID, time.Now().Unix(), ww.Photo.FileType)
	if err := storage.Put(ctx, "wishing-wall", dir, ww.Photo.Byte, ww.Photo.Size, true, ww.Photo.ContentType); err != nil {
		return nil, err
	}

	message := WishingWallMessage{
		ID:      identity.NewID().String(),
		Name:    ww.Name,
		Message: ww.Message,
		Photo:   dir,
		EventID: ww.EventID,
	}

	msgByte, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	err = socketClient.Write(ctx, msgByte)
	if err != nil {
		return nil, err
	}

	return repository.Create(ctx, &message)
}

func (id *EventID) GetAllMessages(ctx context.Context, repository WishingWallMessageRepository) ([]WishingWallMessage, error) {
	return repository.GetAllByEventID(ctx, string(*id))
}
