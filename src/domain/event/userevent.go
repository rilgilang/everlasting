package event

import "time"

type UserEvent struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	EventId   []string  `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
