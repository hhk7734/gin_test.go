package domain

import (
	"time"

	"github.com/hhk7734/gin-test/internal/pkg/validator"
	"github.com/oklog/ulid/v2"
)

func NewUser(emailAddress string) (*User, error) {
	u := &User{
		ID:           ulid.Make(),
		EmailAddress: emailAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

type User struct {
	ID           ulid.ULID
	EmailAddress string `validate:"required,email"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) Validate() error {
	return validator.Validator.Struct(u)
}
