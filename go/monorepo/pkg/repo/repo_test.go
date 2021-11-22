package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name          string
		expectedSpec  Spec
		expectedError string
	}{
		{
			name: "Success",
			expectedSpec: Spec{
				Name: "monorepo",
				Domains: []Domain{
					{
						Name: "auth",
						Subdomains: []Subdomain{
							{Name: "auth"},
						},
					},
					{
						Name: "core",
						Subdomains: []Subdomain{
							{Name: "core"},
						},
					},
				},
				Teams: []Team{
					{Name: "auth-platform"},
					{Name: "core-platform"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			spec, err := Read()

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedSpec, spec)
			} else {
				assert.Empty(t, spec)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
