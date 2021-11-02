package greeting

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-service-horizontal/internal/entity"
)

func TestNewController(t *testing.T) {
	tests := []struct {
		name                    string
		translateGateway        *MockTranslateGateway
		greetingcacheRepository *MockGreetingCacheRepository
		expectedError           string
	}{
		{
			name:                    "OK",
			translateGateway:        &MockTranslateGateway{},
			greetingcacheRepository: &MockGreetingCacheRepository{},
			expectedError:           "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewController(tc.translateGateway, tc.greetingcacheRepository)

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
		name                    string
		translateGateway        *MockTranslateGateway
		greetingcacheRepository *MockGreetingCacheRepository
		ctx                     context.Context
		request                 *entity.GreetRequest
		expectedResponse        *entity.GreetResponse
		expectedError           string
	}{
		{
			name: "Success_FromCache",
			greetingcacheRepository: &MockGreetingCacheRepository{
				LookupMocks: []LookupMock{
					{OutGreeting: "Bonjour"},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &entity.GreetResponse{
				Greeting: "Bonjour, Jane!",
			},
			expectedError: "",
		},
		{
			name: "TranslateFails",
			translateGateway: &MockTranslateGateway{
				TranslateMocks: []TranslateMock{
					{OutError: errors.New("translation error")},
				},
			},
			greetingcacheRepository: &MockGreetingCacheRepository{
				LookupMocks: []LookupMock{
					{OutError: errors.New("not found")},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    "translation error",
		},
		{
			name: "Success_FromAPI",
			translateGateway: &MockTranslateGateway{
				TranslateMocks: []TranslateMock{
					{OutString: "Bonjour"},
				},
			},
			greetingcacheRepository: &MockGreetingCacheRepository{
				LookupMocks: []LookupMock{
					{OutError: errors.New("not found")},
				},
				StoreMocks: []StoreMock{
					{OutError: nil},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &entity.GreetResponse{
				Greeting: "Bonjour, Jane!",
			},
			expectedError: "",
		},
		{
			name: "Success_NoTranlsation",
			translateGateway: &MockTranslateGateway{
				TranslateMocks: []TranslateMock{
					{OutString: ""},
				},
			},
			greetingcacheRepository: &MockGreetingCacheRepository{
				LookupMocks: []LookupMock{
					{OutError: errors.New("not found")},
				},
				StoreMocks: []StoreMock{
					{OutError: nil},
				},
			},
			ctx: context.Background(),
			request: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &entity.GreetResponse{
				Greeting: "Hello, Jane!",
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &controller{
				translateGateway:        tc.translateGateway,
				greetingcacheRepository: tc.greetingcacheRepository,
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
