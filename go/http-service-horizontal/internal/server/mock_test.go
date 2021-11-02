package server

import (
	"context"
	"net/http"
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
	GreetMock struct {
		InResponseWriter http.ResponseWriter
		InRequest        *http.Request
	}

	MockGreetingHandler struct {
		GreetIndex int
		GreetMocks []GreetMock
	}
)

func (m *MockGreetingHandler) Greet(w http.ResponseWriter, r *http.Request) {
	i := m.GreetIndex
	m.GreetIndex++
	m.GreetMocks[i].InResponseWriter = w
	m.GreetMocks[i].InRequest = r
}
