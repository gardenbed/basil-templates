package server

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTP(t *testing.T) {
	tests := []struct {
		name          string
		healthHandler http.Handler
		opts          HTTPOptions
		expectedError string
	}{
		{
			name: "OK",
			healthHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			opts: HTTPOptions{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := NewHTTP(tc.healthHandler, tc.opts)

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

func TestHTTP_String(t *testing.T) {
	s := new(HTTP)
	assert.Equal(t, "http-server", s.String())
}

func TestHTTP_ListenAndServe(t *testing.T) {
	tests := []struct {
		name          string
		s             *HTTP
		expectedError string
	}{
		{
			name: "ListenFails",
			s: &HTTP{
				server: &MockHTTPServer{
					ListenAndServeMocks: []ListenAndServeMock{
						{OutError: errors.New("error on listening")},
					},
				},
			},
			expectedError: "error on listening",
		},
		{
			name: "ServerClosed",
			s: &HTTP{
				server: &MockHTTPServer{
					ListenAndServeMocks: []ListenAndServeMock{
						{OutError: http.ErrServerClosed},
					},
				},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.s.ListenAndServe()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestHTTP_Shutdown(t *testing.T) {
	tests := []struct {
		name          string
		s             *HTTP
		ctx           context.Context
		expectedError string
	}{
		{
			name: "Successful",
			s: &HTTP{
				server: &MockHTTPServer{
					ShutdownMocks: []ShutdownMock{
						{OutError: nil},
					},
				},
			},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name: "Unsuccessful",
			s: &HTTP{
				server: &MockHTTPServer{
					ShutdownMocks: []ShutdownMock{
						{OutError: errors.New("error on shutdown")},
					},
				},
			},
			ctx:           context.Background(),
			expectedError: "error on shutdown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.s.Shutdown(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
