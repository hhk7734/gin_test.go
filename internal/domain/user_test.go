package domain_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hhk7734/gin-test/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserValidate(t *testing.T) {
	cases := []struct {
		User         *domain.User
		InvalidCount int
	}{
		{
			User: &domain.User{
				EmailAddress: "test@test.com",
				FirstName:    "test",
				LastName:     "test",
				PhoneNumber:  "+821012345678",
			},
			InvalidCount: 0,
		},
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
		msg := fmt.Sprintf("case: %d", i)
		err := c.User.Validate()
		ierr := &validator.InvalidValidationError{}
		verr := &validator.ValidationErrors{}
		switch {
		case errors.As(err, verr):
			assert.Equal(t, c.InvalidCount, len(*verr), msg)
		case errors.As(err, &ierr):
			require.Fail(t, "do not allow InvalidValidationError")
		case err == nil:
			if c.InvalidCount != 0 {
				assert.Equal(t, c.InvalidCount, 0, msg)
			}
		case err != nil:
			assert.NoError(t, err, msg)
		}
	}
}
