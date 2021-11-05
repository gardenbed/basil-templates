package github

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name          string
		expectedError string
	}{
		{
			name:          "OK",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := NewService()

			if tc.expectedError == "" {
				assert.NotNil(t, s)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, s)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestService_GetUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.github.com/users/octocat", nil)

	tests := []struct {
		name          string
		client        *MockHTTPClient
		ctx           context.Context
		username      string
		expectedUser  *User
		expectedError string
	}{
		{
			name: "DoFails",
			client: &MockHTTPClient{
				DoMocks: []DoMock{
					{OutError: errors.New("http error")},
				},
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedError: "http error",
		},
		{
			name: "UnexpectedStatusCode",
			client: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							Request:    req,
							StatusCode: 400,
							Body: io.NopCloser(
								strings.NewReader(`{ "error": "invalid request" }`),
							),
						},
					},
				},
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedError: "GET /users/octocat 400: invalid request",
		},
		{
			name: "InvalidResponseBody",
			client: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							Request:    req,
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{`),
							),
						},
					},
				},
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedError: "unexpected EOF",
		},
		{
			name: "Success",
			client: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							Request:    req,
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{ "id": 1, "login": "octocat", "email": "octocat@example.com", "name": "Octocat" }`),
							),
						},
					},
				},
			},
			ctx:      context.Background(),
			username: "octocat",
			expectedUser: &User{
				ID:    1,
				Login: "octocat",
				Email: "octocat@example.com",
				Name:  "Octocat",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Service{
				client: tc.client,
			}

			user, err := s.GetUser(tc.ctx, tc.username)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			} else {
				assert.Empty(t, user)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
