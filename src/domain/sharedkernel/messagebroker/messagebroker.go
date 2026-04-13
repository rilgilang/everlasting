package messagebroker

import "context"

type TaskName string

type MessageBroker interface {
	Produce(ctx context.Context, task TaskName, payload interface{}) error
}
