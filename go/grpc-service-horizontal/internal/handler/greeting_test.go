package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-service-horizontal/internal/entity"
	"grpc-service-horizontal/internal/idl/greetingpb"
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
		ctx                context.Context
		request            *greetingpb.GreetRequest
		expectedResponse   *greetingpb.GreetResponse
		expectedError      string
	}{
		{
			name:             "RequestMappingFails",
			ctx:              context.Background(),
			request:          nil,
			expectedResponse: nil,
			expectedError:    "greet request cannot be nil",
		},
		{
			name: "ControllerFails",
			greetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutError: errors.New("controller error")},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    "controller error",
		},
		{
			name: "ResponseMappingFails",
			greetingController: &MockGreetingController{
				GreetMocks: []GreetMock{
					{OutResponse: nil},
				},
			},
			ctx: context.Background(),
			request: &greetingpb.GreetRequest{
				Name: "Jane",
			},
			expectedResponse: nil,
			expectedError:    "greet response cannot be nil",
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
			handler := &greetingHandler{
				greetingController: tc.greetingController,
			}

			response, err := handler.Greet(tc.ctx, tc.request)

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
