package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gardenbed/basil/httpx"
	"github.com/stretchr/testify/assert"

	"http-service-horizontal/internal/entity"
)

func TestNewGreetingHandler(t *testing.T) {
	tests := []struct {
		name               string
		greetingController *MockGreetingController
		expectedError      string
	}{
		{
			name:               "OK",
			greetingController: &MockGreetingController{},
			expectedError:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, err := NewGreetingHandler(tc.greetingController)

			if tc.expectedError == "" {
				assert.NotNil(t, h)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, h)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGreetingHandler_Greet(t *testing.T) {
	tests := []struct {
		name               string
		greetingController *MockGreetingController
		req                *http.Request
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "RequestDecodingFails",
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{`)),
			expectedStatusCode: 400,
			expectedBody:       "unexpected EOF\n",
		},
		{
			name:               "RequestMappingFails",
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "" }`)),
			expectedStatusCode: 400,
			expectedBody:       "github username cannot be empty\n",
		},
		{
			name: "ControllerFails",
			greetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutError: httpx.NewServerError(errors.New("controller failed"), 500)},
				},
			},
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 500,
			expectedBody:       "controller failed\n",
		},
		{
			name: "ResponseMappingFails",
			greetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutResponse: &entity.GreetResponse{}},
				},
			},
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 500,
			expectedBody:       "greeting cannot be empty\n",
		},
		{
			name: "Success",
			greetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{
						OutResponse: &entity.GreetResponse{
							Greeting: "Hello, Jane!",
						},
					},
				},
			},
			req:                httptest.NewRequest("POST", "/greet", strings.NewReader(`{ "githubUsername": "octocat" }`)),
			expectedStatusCode: 200,
			expectedBody:       "{\"greeting\":\"Hello, Jane!\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := &greetingHandler{
				greetingController: tc.greetingController,
			}

			rec := httptest.NewRecorder()
			handler.Greet(rec, tc.req)

			res := rec.Result()
			b, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			body := string(b)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}
