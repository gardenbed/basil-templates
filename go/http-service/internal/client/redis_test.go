package client

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	c := NewRedis("redis:6379")
	assert.NotNil(t, c)
}

func TestRedis_String(t *testing.T) {
	c := &Redis{}
	str := c.String()

	assert.Equal(t, "redis-client", str)
}

func TestRedis_Connect(t *testing.T) {
	tests := []struct {
		name            string
		universalClient *MockUniversalClient
		expectedError   string
	}{
		{
			name: "PingFails",
			universalClient: &MockUniversalClient{
				PingMocks: []PingMock{
					{OutStatusCmd: redis.NewStatusResult("", errors.New("redis error"))},
				},
			},
			expectedError: "redis error",
		},
		{
			name: "PingSucceeds",
			universalClient: &MockUniversalClient{
				PingMocks: []PingMock{
					{OutStatusCmd: redis.NewStatusResult("", nil)},
				},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Redis{
				UniversalClient: tc.universalClient,
			}

			err := c.Connect()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRedis_Disconnect(t *testing.T) {
	tests := []struct {
		name            string
		universalClient *MockUniversalClient
		ctx             context.Context
		expectedError   string
	}{
		{
			name: "CloseFails",
			universalClient: &MockUniversalClient{
				CloseMocks: []CloseMock{
					{OutError: errors.New("redis error")},
				},
			},
			ctx:           context.Background(),
			expectedError: "redis error",
		},
		{
			name: "CloseSucceeds",
			universalClient: &MockUniversalClient{
				CloseMocks: []CloseMock{
					{OutError: nil},
				},
			},
			ctx:           context.Background(),
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Redis{
				UniversalClient: tc.universalClient,
			}

			err := c.Disconnect(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRedis_HealthCheck(t *testing.T) {
	tests := []struct {
		name            string
		universalClient *MockUniversalClient
		ctx             context.Context
		expectedError   string
	}{
		{
			name: "PingFails",
			universalClient: &MockUniversalClient{
				PingMocks: []PingMock{
					{OutStatusCmd: redis.NewStatusResult("", errors.New("redis error"))},
				},
			},
			ctx:           context.Background(),
			expectedError: "redis error",
		},
		{
			name: "PingSucceeds",
			universalClient: &MockUniversalClient{
				PingMocks: []PingMock{
					{OutStatusCmd: redis.NewStatusResult("", nil)},
				},
			},
			ctx:           context.Background(),
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Redis{
				UniversalClient: tc.universalClient,
			}

			err := c.HealthCheck(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
