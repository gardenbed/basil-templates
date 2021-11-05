package greeting

import (
	"context"

	githubentity "grpc-service-horizontal/internal/entity/github"
)

type (
	ConnectMock struct {
		OutError error
	}

	DisconnectMock struct {
		InContext context.Context
		OutError  error
	}

	MockClient struct {
		StringOut string

		ConnectIndex int
		ConnectMocks []ConnectMock

		DisconnectIndex int
		DisconnectMocks []DisconnectMock
	}
)

func (m *MockClient) String() string {
	return m.StringOut
}

func (m *MockClient) Connect() error {
	i := m.ConnectIndex
	m.ConnectIndex++
	return m.ConnectMocks[i].OutError
}

func (m *MockClient) Disconnect(ctx context.Context) error {
	i := m.DisconnectIndex
	m.DisconnectIndex++
	m.DisconnectMocks[i].InContext = ctx
	return m.DisconnectMocks[i].OutError
}

type (
	HealthCheckMock struct {
		InContext context.Context
		OutError  error
	}

	MockChecker struct {
		StringOut string

		HealthCheckIndex int
		HealthCheckMocks []HealthCheckMock
	}
)

func (m *MockChecker) String() string {
	return m.StringOut
}

func (m *MockChecker) HealthCheck(ctx context.Context) error {
	i := m.HealthCheckIndex
	m.HealthCheckIndex++
	m.HealthCheckMocks[i].InContext = ctx
	return m.HealthCheckMocks[i].OutError
}

type (
	GetUserMock struct {
		InContext  context.Context
		InUsername string
		OutUser    *githubentity.User
		OutError   error
	}

	MockGithubGateway struct {
		MockClient
		MockChecker

		StringOut string

		GetUserIndex int
		GetUserMocks []GetUserMock
	}
)

func (m *MockGithubGateway) String() string {
	return m.StringOut
}

func (m *MockGithubGateway) GetUser(ctx context.Context, username string) (*githubentity.User, error) {
	i := m.GetUserIndex
	m.GetUserIndex++
	m.GetUserMocks[i].InContext = ctx
	m.GetUserMocks[i].InUsername = username
	return m.GetUserMocks[i].OutUser, m.GetUserMocks[i].OutError
}

type (
	StoreMock struct {
		InContext  context.Context
		InUsername string
		InName     string
		OutError   error
	}

	LookupMock struct {
		InContext  context.Context
		InUsername string
		OutName    string
		OutError   error
	}

	MockUserCacheRepository struct {
		MockClient
		MockChecker

		StringOut string

		StoreIndex int
		StoreMocks []StoreMock

		LookupIndex int
		LookupMocks []LookupMock
	}
)

func (m *MockUserCacheRepository) String() string {
	return m.StringOut
}

func (m *MockUserCacheRepository) Store(ctx context.Context, username, name string) error {
	i := m.StoreIndex
	m.StoreIndex++
	m.StoreMocks[i].InContext = ctx
	m.StoreMocks[i].InUsername = username
	m.StoreMocks[i].InName = name
	return m.StoreMocks[i].OutError
}

func (m *MockUserCacheRepository) Lookup(ctx context.Context, username string) (string, error) {
	i := m.LookupIndex
	m.LookupIndex++
	m.LookupMocks[i].InContext = ctx
	m.LookupMocks[i].InUsername = username
	return m.LookupMocks[i].OutName, m.LookupMocks[i].OutError
}
