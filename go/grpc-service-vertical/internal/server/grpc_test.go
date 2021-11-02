package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-service-vertical/internal/idl/greetingpb"
)

func TestNewGRPC(t *testing.T) {
	tests := []struct {
		name            string
		greetingService greetingpb.GreetingServiceServer
		opts            GRPCOptions
		expectedError   string
	}{
		{
			name:            "OK",
			greetingService: &MockGreetingService{},
			opts:            GRPCOptions{},
			expectedError:   "",
		},
		{
			name:            "WithTLS",
			greetingService: &MockGreetingService{},
			opts: GRPCOptions{
				TLSCert:  &tls.Certificate{},
				ClientCA: x509.NewCertPool(),
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := NewGRPC(tc.greetingService, tc.opts)

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

func TestGRPC_String(t *testing.T) {
	s := new(GRPC)
	assert.Equal(t, "grpc-server", s.String())
}

func TestGRPC_ListenAndServe(t *testing.T) {
	tests := []struct {
		name          string
		s             *GRPC
		expectedError string
	}{
		{
			name: "ListenFails",
			s: &GRPC{
				addr: ":-1",
			},
			expectedError: "listen tcp: address -1: invalid port",
		},
		{
			name: "Successful",
			s: &GRPC{
				addr: "127.0.0.1:",
				server: &MockGRPCServer{
					ServeMocks: []ServeMock{
						{OutError: nil},
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

func TestGRPC_Shutdown(t *testing.T) {
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name          string
		s             *GRPC
		ctx           context.Context
		expectedError string
	}{
		{
			name: "Successful",
			s: &GRPC{
				addr:   "127.0.0.1:",
				server: &MockGRPCServer{},
			},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name: "ContextCancelled",
			s: &GRPC{
				addr:   "127.0.0.1:",
				server: &MockGRPCServer{},
			},
			ctx:           cancelledCtx,
			expectedError: "context canceled",
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
