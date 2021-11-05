package usercache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type (
	CloseMock struct {
		OutError error
	}

	PingMock struct {
		InContext    context.Context
		OutStatusCmd *redis.StatusCmd
	}

	GetMock struct {
		InContext    context.Context
		InKey        string
		OutStringCmd *redis.StringCmd
	}

	SetMock struct {
		InContext    context.Context
		InKey        string
		InValue      interface{}
		InExpiration time.Duration
		OutStatusCmd *redis.StatusCmd
	}

	MockRedisClient struct {
		CloseIndex int
		CloseMocks []CloseMock

		PingIndex int
		PingMocks []PingMock

		GetIndex int
		GetMocks []GetMock

		SetIndex int
		SetMocks []SetMock
	}
)

func (m *MockRedisClient) Close() error {
	i := m.CloseIndex
	m.CloseIndex++
	return m.CloseMocks[i].OutError
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	i := m.PingIndex
	m.PingIndex++
	return m.PingMocks[i].OutStatusCmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	i := m.GetIndex
	m.GetIndex++
	return m.GetMocks[i].OutStringCmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	i := m.SetIndex
	m.SetIndex++
	return m.SetMocks[i].OutStatusCmd
}
