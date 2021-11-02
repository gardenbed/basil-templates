package greetingcache

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name          string
		redisAddress  string
		expectedError string
	}{
		{
			name:          "OK",
			redisAddress:  "redis:6379",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewRepository(tc.redisAddress)

			if tc.expectedError == "" {
				assert.NotNil(t, r)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, r)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRepository_String(t *testing.T) {
	r := new(repository)
	assert.Equal(t, "greetingcache-repository", r.String())
}

func TestRepository_Connect(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockRedisClient
		expectedError string
	}{
		{
			name: "PingFails",
			client: &MockRedisClient{
				PingMocks: []PingMock{
					{OutStatusCmd: redis.NewStatusResult("", errors.New("redis error"))},
				},
			},
			expectedError: "redis error",
		},
		{
			name: "Success",
			client: &MockRedisClient{
				PingMocks: []PingMock{
					{OutStatusCmd: redis.NewStatusResult("", nil)},
				},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &repository{
				client: tc.client,
			}

			err := r.Connect()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRepository_Disconnect(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockRedisClient
		ctx           context.Context
		expectedError string
	}{
		{
			name: "CloseFails",
			client: &MockRedisClient{
				CloseMocks: []CloseMock{
					{OutError: errors.New("redis error")},
				},
			},
			ctx:           context.Background(),
			expectedError: "redis error",
		},
		{
			name: "Success",
			client: &MockRedisClient{
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
			r := &repository{
				client: tc.client,
			}

			err := r.Disconnect(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRepository_HealthCheck(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockRedisClient
		ctx           context.Context
		expectedError string
	}{
		{
			name: "PingFails",
			client: &MockRedisClient{
				PingMocks: []PingMock{
					{OutStatusCmd: redis.NewStatusResult("", errors.New("redis error"))},
				},
			},
			ctx:           context.Background(),
			expectedError: "redis error",
		},
		{
			name: "Success",
			client: &MockRedisClient{
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
			r := &repository{
				client: tc.client,
			}

			err := r.HealthCheck(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRepository_Store(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockRedisClient
		ctx           context.Context
		lang          string
		greeting      string
		expectedError string
	}{
		{
			name:          "NoLanguage",
			client:        &MockRedisClient{},
			ctx:           context.Background(),
			lang:          "",
			greeting:      "",
			expectedError: "no language code",
		},
		{
			name:          "NoGreeting",
			client:        &MockRedisClient{},
			ctx:           context.Background(),
			lang:          "fr",
			greeting:      "",
			expectedError: "no greeting value",
		},
		{
			name: "SetFails",
			client: &MockRedisClient{
				SetMocks: []SetMock{
					{OutStatusCmd: redis.NewStatusResult("", errors.New("redis error"))},
				},
			},
			ctx:           context.Background(),
			lang:          "fr",
			greeting:      "Hello, World!",
			expectedError: "redis error",
		},
		{
			name: "Success",
			client: &MockRedisClient{
				SetMocks: []SetMock{
					{OutStatusCmd: redis.NewStatusResult("", nil)},
				},
			},
			ctx:           context.Background(),
			lang:          "fr",
			greeting:      "Hello, World!",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &repository{
				client: tc.client,
			}

			err := r.Store(tc.ctx, tc.lang, tc.greeting)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRepository_Lookup(t *testing.T) {
	tests := []struct {
		name             string
		client           *MockRedisClient
		ctx              context.Context
		lang             string
		expectedGreeting string
		expectedError    string
	}{
		{
			name:             "NoLanguage",
			client:           &MockRedisClient{},
			ctx:              context.Background(),
			lang:             "",
			expectedGreeting: "",
			expectedError:    "no language code",
		},
		{
			name: "GetFails",
			client: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("", errors.New("redis error"))},
				},
			},
			ctx:              context.Background(),
			lang:             "fr",
			expectedGreeting: "",
			expectedError:    "redis error",
		},
		{
			name: "Success",
			client: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("Bonjour", nil)},
				},
			},
			ctx:              context.Background(),
			lang:             "fr",
			expectedGreeting: "Bonjour",
			expectedError:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &repository{
				client: tc.client,
			}

			greeting, err := r.Lookup(tc.ctx, tc.lang)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedGreeting, greeting)
			} else {
				assert.Empty(t, greeting)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
