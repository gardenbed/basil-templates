package greet

import (
	"context"

	"command-line-app/internal/github"
)

type (
	GetUserMock struct {
		InContext  context.Context
		InUsername string
		OutUser    *github.User
		OutError   error
	}

	MockGithubService struct {
		GetUserIndex int
		GetUserMocks []GetUserMock
	}
)

func (m *MockGithubService) GetUser(ctx context.Context, username string) (*github.User, error) {
	i := m.GetUserIndex
	m.GetUserIndex++
	m.GetUserMocks[i].InContext = ctx
	m.GetUserMocks[i].InUsername = username
	return m.GetUserMocks[i].OutUser, m.GetUserMocks[i].OutError
}
