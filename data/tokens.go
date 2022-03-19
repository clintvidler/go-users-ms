package data

import (
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type Token struct {
	ID        uint      `validate:"required"`
	UserEmail string    `validate:"required,email"`
	Token     string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	ExpiredAt time.Time `validate:"required"`
}

func (t *Token) Validate() error {
	validate := validator.New()

	return validate.Struct(t)
}

func (s *Store) ReadOneToken(email, token string) (t Token, err error) {
	s.db.Where("user_email = ? and token = ? and expired_at >= now()", email, token).Last(&t)

	err = t.Validate()

	log.Println(err)

	return
}

func (s *Store) DeleteOneToken(email string) (err error) {
	if email == "" {
		err = errors.New("email must not be blank")
		return
	}

	s.db.Delete(Token{}, "user_email", email)

	return
}
