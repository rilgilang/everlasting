package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	messagebrokerDomain "everlasting/src/domain/sharedkernel/messagebroker"
	"everlasting/src/infrastructure/pkg"

	amqpDriver "github.com/rabbitmq/amqp091-go"

	"everlasting/src/infrastructure/pkg/logger"
)

// Event definition
type (
	Event interface {
		Handle(ctx context.Context, payload []byte) (err error)
	}

	Events map[messagebrokerDomain.TaskName]Event
)

func (es Events) Handle(ctx context.Context, name messagebrokerDomain.TaskName, payload []byte) (err error) {
	if x, found := es[name]; found {
		return x.Handle(ctx, payload)
	}
	return err
}

// Message broker definition
type (
	Body struct {
		TaskName messagebrokerDomain.TaskName
		Payload  interface{}
	}

	MessageBroker struct {
		config *pkg.Config
		logger *logger.AppLogger
	}
)

var DEFAULT_QUEUE_NAME = "everlasting_queue"
var DEFAULT_EXCHANGE_NAME = "everlasting_exchange"
var DEFAULT_EXCHANGE_TYPE = "x-delayed-message"
var DEFAULT_EXCHANGE_ARGS = amqpDriver.Table{
	"x-delayed-type": "direct",
}

func NewMessageBroker(config *pkg.Config, logger *logger.AppLogger) (result *MessageBroker) {
	return &MessageBroker{config, logger}
}

func (mb *MessageBroker) connect(ctx context.Context) (conn *amqpDriver.Connection, channel *amqpDriver.Channel, err error) {
	conn, err = amqpDriver.Dial(mb.config.GenerateAMQPConnectionString())
	if err != nil {
		mb.logger.Error(ctx, "message_broker", err.Error())
		return conn, channel, err
	}

	channel, err = conn.Channel()
	if err != nil {
		mb.logger.Error(ctx, "message_broker", err.Error())
		return conn, channel, err
	}

	// Define exchange
	err = channel.ExchangeDeclare(
		DEFAULT_EXCHANGE_NAME, // name
		DEFAULT_EXCHANGE_TYPE, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		DEFAULT_EXCHANGE_ARGS, // arguments
	)
	if err != nil {
		mb.logger.Error(ctx, "message_broker", err.Error())
		return conn, channel, err
	}

	// Define queue
	queueName := DEFAULT_QUEUE_NAME
	if mb.config.RabbitMQQueueName != "" {
		queueName = mb.config.RabbitMQQueueName
	}

	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		mb.logger.Error(ctx, "message_broker", err.Error())
		return conn, channel, err
	}

	err = channel.QueueBind(
		queue.Name,            // queue name
		"",                    // routing key
		DEFAULT_EXCHANGE_NAME, // exchange
		false,
		nil)

	if err != nil {
		mb.logger.Error(ctx, "message_broker", err.Error())
		return conn, channel, err
	}

	return conn, channel, err

}

func (mb *MessageBroker) Produce(ctx context.Context, task messagebrokerDomain.TaskName, payload interface{}) (err error) {
	body, err := json.Marshal(Body{
		TaskName: task,
		Payload:  payload,
	})
	if err != nil {
		mb.logger.Error(ctx, "message_broker", err.Error())
		return err
	}

	conn, channel, err := mb.connect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
		channel.Close()
	}()

	err = channel.PublishWithContext(ctx,
		DEFAULT_EXCHANGE_NAME,
		"",
		false,
		false,
		amqpDriver.Publishing{
			DeliveryMode: amqpDriver.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)

	if err != nil {
		mb.logger.Error(ctx, "message_broker", err.Error())
		return err
	}

	return err
}

func (mb *MessageBroker) Consume(ev Events) (err error) {
	ctx := context.Background()

	conn, channel, err := mb.connect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
		channel.Close()
	}()

	msgs, err := channel.Consume(
		mb.config.RabbitMQQueueName, // queue
		"aicare_klondike",           // consumer
		true,                        // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)

	if err != nil {
		mb.logger.Fatal(ctx, "message_broker", err.Error())
		return err
	}

	var forever chan struct{}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				mb.logger.Error(ctx, "amqp_consumer", fmt.Sprintf("recovered: %v", err))
			}
		}()
		for d := range msgs {
			var body Body
			ctx := logger.NewAppLoggerContextFromAMQPEvent(body.TaskName).GetContext()
			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				mb.logger.Error(ctx, "message_broker", err.Error())
				continue
			}

			payload, err := json.Marshal(body.Payload)
			if err != nil {
				mb.logger.Error(ctx, "message_broker", err.Error())
				continue
			}

			ev.Handle(ctx, body.TaskName, payload)
		}
	}()

	<-forever
	return err
}
