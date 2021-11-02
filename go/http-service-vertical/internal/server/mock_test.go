package server

import (
	"context"
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
