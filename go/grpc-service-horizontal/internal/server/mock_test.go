package server

import (
	"context"
	"net"

	"grpc-service-horizontal/internal/idl/greetingpb"
)

type (
	ListenAndServeMock struct {
		OutError error
	}

	ShutdownMock struct {
		InContext context.Context
		OutError  error
	}

	MockHTTPServer struct {
		ListenAndServeIndex int
		ListenAndServeMocks []ListenAndServeMock

		ShutdownIndex int
		ShutdownMocks []ShutdownMock
	}
)

func (m *MockHTTPServer) ListenAndServe() error {
	i := m.ListenAndServeIndex
	m.ListenAndServeIndex++
	return m.ListenAndServeMocks[i].OutError
}

func (m *MockHTTPServer) Shutdown(ctx context.Context) error {
	i := m.ShutdownIndex
	m.ShutdownIndex++
	m.ShutdownMocks[i].InContext = ctx
	return m.ShutdownMocks[i].OutError
}

type (
	ServeMock struct {
		InListener net.Listener
		OutError   error
	}

	GracefulStopMock struct {
	}

	MockGRPCServer struct {
		GracefulStopIndex int
		GracefulStopMocks []GracefulStopMock

		ServeIndex int
		ServeMocks []ServeMock
	}
)

func (m *MockGRPCServer) GracefulStop() {
	m.GracefulStopIndex++
}

func (m *MockGRPCServer) Serve(lis net.Listener) error {
	i := m.ServeIndex
	m.ServeIndex++
	m.ServeMocks[i].InListener = lis
	return m.ServeMocks[i].OutError
}

type (
	GreetMock struct {
		InContext   context.Context
		InRequest   *greetingpb.GreetRequest
		OutResponse *greetingpb.GreetResponse
		OutError    error
	}

	MockGreetingHandler struct {
		GreetIndex int
		GreetMocks []GreetMock
	}
)

func (m *MockGreetingHandler) Greet(ctx context.Context, req *greetingpb.GreetRequest) (*greetingpb.GreetResponse, error) {
	i := m.GreetIndex
	m.GreetIndex++
	m.GreetMocks[i].InContext = ctx
	m.GreetMocks[i].InRequest = req
	return m.GreetMocks[i].OutResponse, m.GreetMocks[i].OutError
}
