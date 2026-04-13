package validator

import (
	"unicode"

	"everlasting/src/domain/user"

	"github.com/go-playground/validator"
)

var PasswordCustomValidator = func(fl validator.FieldLevel) bool {
	var hasLower, hasUpper, hasNumber, hasSpecial bool
	for _, char := range fl.Field().Interface().(user.Password) {
		if unicode.IsLetter(char) {
			if unicode.IsLower(char) {
				hasLower = true
			} else if unicode.IsUpper(char) {
				hasUpper = true
			}
		} else if unicode.IsNumber(char) {
			hasNumber = true
		} else {
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasNumber && hasSpecial
}
