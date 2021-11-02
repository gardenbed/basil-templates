package greeting

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"grpc-service-vertical/internal/idl/greetingpb"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name          string
		httpClient    *MockHTTPClient
		redisClient   *MockRedisClient
		expectedError string
	}{
		{
			name:          "OK",
			httpClient:    &MockHTTPClient{},
			redisClient:   &MockRedisClient{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := NewService(tc.httpClient, tc.redisClient)

			if tc.expectedError == "" {
				assert.NotNil(t, s)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, s)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestService_Greet(t *testing.T) {
	tests := []struct {
		name             string
		httpClient       *MockHTTPClient
		redisClient      *MockRedisClient
		ctx              context.Context
		request          *greetingpb.GreetRequest
		expectedResponse *greetingpb.GreetResponse
		expectedError    string
	}{
		{
			name: "Success_FromCache",
			redisClient: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("Bonjour", nil)},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &greetingpb.GreetResponse{
				Greeting: "Bonjour, Jane!",
			},
			expectedError: "",
		},
		{
			name: "HTTPCallFails",
			httpClient: &MockHTTPClient{
				DoMocks: []DoMock{
					{OutError: errors.New("http error")},
				},
			},
			redisClient: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("", errors.New("redis error"))},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    "http error",
		},
		{
			name: "InvalidResponseBody",
			httpClient: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{`),
							),
						},
					},
				},
			},
			redisClient: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("", errors.New("redis error"))},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    "unexpected EOF",
		},
		{
			name: "Success_FromAPI",
			httpClient: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{ "translatedText": "Bonjour" }`),
							),
						},
					},
				},
			},
			redisClient: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("", errors.New("redis error"))},
				},
				SetMocks: []SetMock{
					{OutStatusCmd: redis.NewStatusResult("", nil)},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &greetingpb.GreetResponse{
				Greeting: "Bonjour, Jane!",
			},
			expectedError: "",
		},
		{
			name: "Success_NoTranlsation",
			httpClient: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{}`),
							),
						},
					},
				},
			},
			redisClient: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("", errors.New("redis error"))},
				},
				SetMocks: []SetMock{
					{OutStatusCmd: redis.NewStatusResult("", nil)},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: &greetingpb.GreetResponse{
				Greeting: "Hello, Jane!",
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &service{
				httpClient:  tc.httpClient,
				redisClient: tc.redisClient,
			}

			response, err := s.Greet(tc.ctx, tc.request)

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
