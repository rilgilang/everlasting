package resetpassword

import (
	"context"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/sharedkernel/messagebroker"
	"everlasting/src/domain/user"
)

const TaskSendResetPasswordRequest messagebroker.TaskName = "send_reset_pasword_request"

type ResetPasswordRequest struct {
	Email user.Email `json:"email" example:"xxx@yyy.zzz" validate:"required,email"`
}

func NewResetPasswordRequest(email user.Email) *ResetPasswordRequest {
	return &ResetPasswordRequest{
		Email: email,
	}
}

func (rpr *ResetPasswordRequest) IsHasMatchedUserIn(ctx context.Context, userRepo user.UserRepository) (result bool, err error) {
	_, err = rpr.Email.GetMatchedUserIn(ctx, userRepo)
	if err != nil {
		// If there is no matched email in data storage. Expected persistence should return errDomain.NotFoundEntityError
		if err == errDomain.ErrUserNotFound {
			return result, nil
		}
		// Otherwise ...
		return result, err
	}

	return true, err
}

// Following is method to produce change password instruction to message broker
// Message will be consumed by message broker consumer
func (rpr *ResetPasswordRequest) PutInstructionQueueIn(ctx context.Context, broker messagebroker.MessageBroker) (err error) {
	return broker.Produce(ctx, TaskSendResetPasswordRequest, rpr)
}
