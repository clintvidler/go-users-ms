package data

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" gorm:"primaryKey;<-:create" validate:"required,email"`
	Password     []byte `json:"-" validate:"required"`
	FailedLogins int
	BlockedUntil time.Time
}

func (u *User) Validate() error {
	validate := validator.New()

	return validate.Struct(u)
}

func (u *User) SetPassword(password string) (err error) {
	if password == "" {
		err = errors.New("password is required")
		return
	}
	u.Password, err = bcrypt.GenerateFromPassword([]byte(password), 12)
	return
}

func (u *User) ComparePassword(password string) (err error) {
	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	return
}

func (s *Store) CreateOneUser(u *User) (err error) {
	err = s.db.Create(&u).Error
	return
}

func (s *Store) ReadOneUser(email string) (u User, err error) {
	s.db.Where("email = ?", email).First(&u)

	err = u.Validate()
	return
}

func (s *Store) UpdateOneUser(u User) (err error) {
	err = u.Validate()

	s.db.Save(&u)

	return
}

func (s *Store) Login(u User, token string, expires time.Time) (err error) {
	err = u.Validate()

	if err != nil {
		return
	}

	if token == "" {
		return errors.New("token is required")
	}

	if expires.Before(time.Now()) {
		return errors.New("expires must be in the future")
	}

	t := Token{
		UserEmail: u.Email,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiredAt: expires,
	}

	s.db.Create(&t)

	return
}

func (s *Store) Logout(u User) (err error) {
	err = u.Validate()

	if err != nil {
		return
	}

	_, err = s.ReadOneUser(u.Email)

	if err != nil {
		return
	}

	s.db.Delete(Token{}, "user_email", &u.Email)

	return
}
