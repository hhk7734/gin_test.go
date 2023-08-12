package domain_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hhk7734/gin-test/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserShouldValid(t *testing.T) {
	// Given
	cases := []struct {
		User *domain.User
	}{
		{
			User: &domain.User{
				EmailAddress: "test@test.com",
				FirstName:    "test",
				LastName:     "test",
				PhoneNumber:  "+821012345678",
			},
		},
	}

	for i, c := range cases {
		// When
		err := c.User.Validate()

		// Then
		msg := fmt.Sprintf("case: %d", i)
		assert.NoError(t, err, msg)
	}
}

func TestInvalidUserShouldReturnOnlyValidationErrors(t *testing.T) {
	// Given
	cases := []struct {
		User         *domain.User
		InvalidCount int
	}{
		{
			User:         &domain.User{},
			InvalidCount: 3,
		},
		{
			User: &domain.User{
				EmailAddress: "test",
				FirstName:    strings.Repeat("a", 256),
				LastName:     strings.Repeat("a", 256),
				PhoneNumber:  "010-1234-5678",
			},
			InvalidCount: 4,
		},
	}

	for i, c := range cases {
		// When
		err := c.User.Validate()

		// Then
		msg := fmt.Sprintf("case: %d", i)
		verr := validator.ValidationErrors{}
		if assert.True(t, errors.As(err, &verr), msg) {
			assert.Len(t, verr, c.InvalidCount, msg)
		}
	}
}
