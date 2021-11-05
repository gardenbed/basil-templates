package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	tests := []struct {
		name           string
		entity         User
		expectedString string
	}{
		{
			name: "OK",
			entity: User{
				ID:    1,
				Login: "octocat",
				Email: "octocat@example.com",
				Name:  "Octocat",
			},
			expectedString: "User{id=1 login=octocat email=octocat@example.com name=Octocat}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.entity.String())
		})
	}
}
