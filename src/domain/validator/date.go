package validator

import (
	"time"

	"github.com/go-playground/validator"
)

var DateValidator = func(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	if value == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}
