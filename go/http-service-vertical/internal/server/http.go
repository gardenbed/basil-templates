package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gardenbed/basil/httpx"
	"github.com/gorilla/mux"

	"http-service-vertical/internal/service/greeting"
)

const defaultHTTPPort = 8080

// httpServer is an interface for http.Server struct.
type httpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HTTP is an http server implementing the graceful.Server interface.
type HTTP struct {
	server httpServer
}

// HTTPOptions are optional settings for creating an http server.
type HTTPOptions struct {
	// The port number for the HTTP server.
	// The default port number is 8080.
	Port uint16
	// HTTP middleware for handlers.
	Middleware []httpx.Middleware
}

// NewHTTP creates a new http Server.
func NewHTTP(healthHandler http.Handler, greetingService *greeting.Service, opts HTTPOptions) (*HTTP, error) {
	if opts.Port == 0 {
		opts.Port = defaultHTTPPort
	}

	router := mux.NewRouter()
	router.Path("/health").Handler(healthHandler)
	greetingService.RegisterRoutes(router, opts.Middleware...)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: router,
	}

	return &HTTP{
		server: server,
	}, nil
}

// String returns the name of the server.
func (s *HTTP) String() string {
	return "http-server"
}

// ListenAndServe starts listening for incoming requests synchronously.
// It blocks the current goroutine until an error is returned.
func (s *HTTP) ListenAndServe() error {
	// Synchronous/Blocking
	// ListenAndServe always returns a non-nil error
	// After Shutdown or Close, the returned error is ErrServerClosed
	err := s.server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server.
// It stops accepting new conenctions and blocks the current goroutine until all the pending requests are completed.
// If the context is cancelled, an error will be returned.
func (s *HTTP) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
