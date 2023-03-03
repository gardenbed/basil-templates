package client

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type (
	CloseMock struct {
		OutError error
	}

	PingMock struct {
		InContext    context.Context
		OutStatusCmd *redis.StatusCmd
	}

	MockUniversalClient struct {
		redis.UniversalClient

		CloseIndex int
		CloseMocks []CloseMock

		PingIndex int
		PingMocks []PingMock
	}
)

func (m *MockUniversalClient) Close() error {
	i := m.CloseIndex
	m.CloseIndex++
	return m.CloseMocks[i].OutError
}

func (m *MockUniversalClient) Ping(ctx context.Context) *redis.StatusCmd {
	i := m.PingIndex
	m.PingIndex++
	return m.PingMocks[i].OutStatusCmd
}
