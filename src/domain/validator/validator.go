package validator

import (
	errDomain "everlasting/src/domain/error"

	"github.com/go-playground/validator"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()

	// register your custom validation here
	validate.RegisterValidation("date", DateValidator)
	validate.RegisterValidation("required_string", RequiredStringValidator)
	validate.RegisterValidation("password_custom_validator", PasswordCustomValidator)

	return &CustomValidator{
		validator: validate,
	}
}

func (v *CustomValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		if ve, ok := err.((validator.ValidationErrors)); ok {
			return ve
		}
		return errDomain.ErrValidation
	}
	return nil
}
