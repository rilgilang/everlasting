package resetpassword

import (
	"context"

	"everlasting/src/domain/sharedkernel/datatype"
	"everlasting/src/domain/user"
	userDomain "everlasting/src/domain/user"
)

const UIDContextKey datatype.ContextKey = "UID"

type SubmitPasswordRequest struct {
	Password        userDomain.Password `json:"password" example:"1234qweR!" validate:"required,password_custom_validator"`
	PasswordConfirm string              `json:"password_confirm" example:"1234qweR!" validate:"eqfield=Password"`
}

func (s *SubmitPasswordRequest) SaveTo(ctx context.Context, userRepo userDomain.UserRepository) (err error) {
	uid := ctx.Value(UIDContextKey).(user.UserID)

	user, err := uid.GetDetailFrom(ctx, userRepo)
	if err != nil {
		return err
	}

	err = user.SetPassword(s.Password)
	if err != nil {
		return err
	}

	_, err = user.UpdateTo(ctx, userRepo)
	return err
}
