package domain_test

import (
	"testing"

	"github.com/hhk7734/gin-test/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {
	cases := []struct {
		user    *domain.User
		isValid bool
	}{
		{
			user: &domain.User{
				EmailAddress: "",
			},
			isValid: false,
		},
		{
			user: &domain.User{
				EmailAddress: "test",
			},
			isValid: false,
		},
		{
			user: &domain.User{
				EmailAddress: "test@test.com",
			},
			isValid: true,
		},
	}

	for _, c := range cases {
		err := c.user.Validate()
		if c.isValid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
