package validator

import (
	"strings"

	"github.com/go-playground/validator"
)

var RequiredStringValidator = func(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	return strings.Trim(value, " ") != ""
}
