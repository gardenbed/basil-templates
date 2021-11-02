package greeting

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type (
	DoMock struct {
		InRequest   *http.Request
		OutResponse *http.Response
		OutError    error
	}

	MockHTTPClient struct {
		DoIndex int
		DoMocks []DoMock
	}
)

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	i := m.DoIndex
	m.DoIndex++
	m.DoMocks[i].InRequest = req
	return m.DoMocks[i].OutResponse, m.DoMocks[i].OutError
}

type (
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
		GetIndex int
		GetMocks []GetMock

		SetIndex int
		SetMocks []SetMock
	}
)

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

type (
	MockMiddleware struct{}
)

func (m *MockMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return next
}
