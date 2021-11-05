package greeting

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-service-horizontal/internal/entity"
	githubentity "grpc-service-horizontal/internal/entity/github"
)

func TestNewController(t *testing.T) {
	tests := []struct {
		name                string
		githubGateway       *MockGithubGateway
		usercacheRepository *MockUserCacheRepository
		expectedError       string
	}{
		{
			name:                "OK",
			githubGateway:       &MockGithubGateway{},
			usercacheRepository: &MockUserCacheRepository{},
			expectedError:       "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewController(tc.githubGateway, tc.usercacheRepository)

			if tc.expectedError == "" {
				assert.NotNil(t, c)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, c)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestController_Greet(t *testing.T) {
	tests := []struct {
		name                string
		githubGateway       *MockGithubGateway
		usercacheRepository *MockUserCacheRepository
		ctx                 context.Context
		request             *entity.GreetRequest
		expectedResponse    *entity.GreetResponse
		expectedError       string
	}{
		{
			name: "Success_FromCache",
			usercacheRepository: &MockUserCacheRepository{
				LookupMocks: []LookupMock{
					{OutName: "Octocat"},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				GithubUsername: "octocat",
			},
			expectedResponse: &entity.GreetResponse{
				Greeting: "Hello, Octocat!",
			},
			expectedError: "",
		},
		{
			name: "GetUserFails",
			githubGateway: &MockGithubGateway{
				GetUserMocks: []GetUserMock{
					{OutError: errors.New("github error")},
				},
			},
			usercacheRepository: &MockUserCacheRepository{
				LookupMocks: []LookupMock{
					{OutError: errors.New("not found")},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				GithubUsername: "octocat",
			},
			expectedResponse: nil,
			expectedError:    "github error",
		},
		{
			name: "Success_FromAPI",
			githubGateway: &MockGithubGateway{
				GetUserMocks: []GetUserMock{{
					OutUser: &githubentity.User{
						ID:    1,
						Login: "octocat",
						Email: "octocat@example.com",
						Name:  "Octocat",
					},
				},
				},
			},
			usercacheRepository: &MockUserCacheRepository{
				LookupMocks: []LookupMock{
					{OutError: errors.New("not found")},
				},
				StoreMocks: []StoreMock{
					{OutError: nil},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				GithubUsername: "octocat",
			},
			expectedResponse: &entity.GreetResponse{
				Greeting: "Hello, Octocat!",
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &controller{
				githubGateway:       tc.githubGateway,
				usercacheRepository: tc.usercacheRepository,
			}

			response, err := c.Greet(tc.ctx, tc.request)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			} else {
				assert.Nil(t, response)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
