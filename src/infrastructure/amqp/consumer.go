package amqp

import (
	wishingWallEvent "everlasting/src/domain/event"
	messagebrokerDomain "everlasting/src/domain/sharedkernel/messagebroker"
	"everlasting/src/domain/user/resetpassword"
	"everlasting/src/infrastructure/amqp/event"
	"everlasting/src/infrastructure/pkg"

	"github.com/sarulabs/di"
)

func Consume(container di.Container, config *pkg.Config) (err error) {
	events := map[messagebrokerDomain.TaskName]Event{
		resetpassword.TaskSendResetPasswordRequest: event.NewResetPassword(container, config),
		wishingWallEvent.WishingWallMessageTask:    event.NewWishingWallMessage(container, config),
	}

	err = container.Get("pkg.messagebroker.amqp").(*MessageBroker).Consume(events)

	if err != nil {
		return err
	}
	return err
}
