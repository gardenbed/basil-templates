package github

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	githubentity "grpc-service-horizontal/internal/entity/github"
)

func TestNewGateway(t *testing.T) {
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
			g, err := NewGateway()

			if tc.expectedError == "" {
				assert.NotNil(t, g)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, g)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_String(t *testing.T) {
	g := new(gateway)
	assert.Equal(t, "github-gateway", g.String())
}

func TestGateway_Connect(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockHTTPClient
		expectedError string
	}{
		{
			name:          "OK",
			client:        &MockHTTPClient{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			err := g.Connect()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_Disconnect(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockHTTPClient
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "OK",
			client:        &MockHTTPClient{},
			ctx:           context.Background(),
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			err := g.Disconnect(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_HealthCheck(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockHTTPClient
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "OK",
			client:        &MockHTTPClient{},
			ctx:           context.Background(),
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			err := g.HealthCheck(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_GetUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.github.com/users/octocat", nil)

	tests := []struct {
		name          string
		client        *MockHTTPClient
		ctx           context.Context
		username      string
		expectedUser  *githubentity.User
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
			name: "InvalidResponseStatusCode",
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
			expectedUser: &githubentity.User{
				ID:    1,
				Login: "octocat",
				Email: "octocat@example.com",
				Name:  "Octocat",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			text, err := g.GetUser(tc.ctx, tc.username)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, text)
			} else {
				assert.Empty(t, text)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
