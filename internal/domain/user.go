package domain

import (
	"time"

	"github.com/hhk7734/gin_test.go/internal/pkg/validator"
	"github.com/oklog/ulid/v2"
)

type User struct {
	ID           ulid.ULID
	EmailAddress string `validate:"required,email,max=255"`
	FirstName    string `validate:"required,max=255"`
	LastName     string `validate:"required,max=255"`
	PhoneNumber  string `validate:"omitempty,e164"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(emailAddress string) *User {
	return &User{
		ID:           ulid.Make(),
		EmailAddress: emailAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (u *User) Validate() error {
	return validator.Validator.Struct(u)
}
