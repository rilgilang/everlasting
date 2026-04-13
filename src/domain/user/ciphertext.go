package user

import (
	"golang.org/x/crypto/bcrypt"
)

// CipherText Definition (Hashed version of password)
type CipherText string

func NewCipherTextFromPassword(p Password) (CipherText, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 1)
	return CipherText(bytes), err
}

func (h CipherText) VerifyPassword(p Password) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
}
