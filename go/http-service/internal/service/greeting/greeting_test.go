package greeting

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gardenbed/basil/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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

func TestService_RegisterRoutes(t *testing.T) {
	tests := []struct {
		name       string
		router     *mux.Router
		middleware []httpx.Middleware
	}{
		{
			name:   "OK",
			router: mux.NewRouter(),
			middleware: []httpx.Middleware{
				&MockMiddleware{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Service{}

			s.RegisterRoutes(tc.router, tc.middleware...)

			assert.NotEmpty(t, tc.router.Get("Greet"))
		})
	}
}

func TestService_Greet(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.github.com/users/octocat", nil)

	tests := []struct {
		name               string
		httpClient         *MockHTTPClient
		redisClient        *MockRedisClient
		ctx                context.Context
		r                  *http.Request
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "InvalidRequest",
			ctx:                context.Background(),
			r:                  httptest.NewRequest("POST", "/greet", strings.NewReader(`{`)),
			expectedStatusCode: 400,
			expectedBody:       "unexpected EOF\n",
		},
		{
			name: "Success_FromCache",
			redisClient: &MockRedisClient{
				GetMocks: []GetMock{
					{OutStringCmd: redis.NewStringResult("Octocat", nil)},
				},
			},
			ctx:                context.Background(),
			r:                  httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 200,
			expectedBody:       "{\"greeting\":\"Hello, Octocat!\"}\n",
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
			ctx:                context.Background(),
			r:                  httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 500,
			expectedBody:       "http error\n",
		},
		{
			name: "UnexpectedStatusCode",
			httpClient: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							Request:    req,
							StatusCode: 400,
							Body: io.NopCloser(
								strings.NewReader(`{ "error": "invalid request" }`),
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
			ctx:                context.Background(),
			r:                  httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 500,
			expectedBody:       "GET /users/octocat 400: invalid request\n",
		},
		{
			name: "InvalidResponseBody",
			httpClient: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							Request:    req,
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
			ctx:                context.Background(),
			r:                  httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 500,
			expectedBody:       "unexpected EOF\n",
		},
		{
			name: "Success_FromAPI",
			httpClient: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							Request:    req,
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{ "id": 1, "login": "octocat", "email": "octocat@example.com", "name": "Octocat" }`),
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
			ctx:                context.Background(),
			r:                  httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 200,
			expectedBody:       "{\"greeting\":\"Hello, Octocat!\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Service{
				httpClient:  tc.httpClient,
				redisClient: tc.redisClient,
			}

			rec := httptest.NewRecorder()
			s.Greet(rec, tc.r)

			res := rec.Result()
			b, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			body := string(b)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}
