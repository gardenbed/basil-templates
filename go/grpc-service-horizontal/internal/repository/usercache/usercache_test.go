package usercache

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
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
	assert.Equal(t, "usercache-repository", r.String())
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
		testname      string
		client        *MockRedisClient
		ctx           context.Context
		username      string
		name          string
		expectedError string
	}{
		{
			testname:      "NoUsername",
			client:        &MockRedisClient{},
			ctx:           context.Background(),
			username:      "",
			name:          "",
			expectedError: "no username",
		},
		{
			testname:      "NoName",
			client:        &MockRedisClient{},
			ctx:           context.Background(),
			username:      "octocat",
			name:          "",
			expectedError: "no name",
		},
		{
			testname: "SetFails",
			client: &MockRedisClient{
				SetMocks: []SetMock{
					{OutStatusCmd: redis.NewStatusResult("", errors.New("redis error"))},
				},
			},
			ctx:           context.Background(),
			username:      "octocat",
			name:          "Octocat",
			expectedError: "redis error",
		},
		{
			testname: "Success",
			client: &MockRedisClient{
				SetMocks: []SetMock{
					{OutStatusCmd: redis.NewStatusResult("", nil)},
				},
			},
			ctx:           context.Background(),
			username:      "octocat",
			name:          "Octocat",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			r := &repository{
				client: tc.client,
			}

			err := r.Store(tc.ctx, tc.username, tc.name)

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
		testname      string
		client        *MockRedisClient
		ctx           context.Context
		username      string
		expectedName  string
		expectedError string
	}{
		{
			testname:      "NoUsername",
			client:        &MockRedisClient{},
			ctx:           context.Background(),
			username:      "",
			expectedName:  "",
			expectedError: "no username",
		},
		{
			testname: "GetFails",
			client: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("", errors.New("redis error"))},
				},
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedName:  "",
			expectedError: "redis error",
		},
		{
			testname: "Success",
			client: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("Octocat", nil)},
				},
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedName:  "Octocat",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			r := &repository{
				client: tc.client,
			}

			name, err := r.Lookup(tc.ctx, tc.username)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedName, name)
			} else {
				assert.Empty(t, name)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
