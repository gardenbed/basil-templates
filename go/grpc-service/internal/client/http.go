package client

import (
	"context"
	"net/http"
	"time"
)

// HTTP is an http.Client that implements the graceful.Client and graceful.Client health.Checker interfaces.
type HTTP struct {
	*http.Client
}

// NewHTTP creates a new http client.
func NewHTTP() *HTTP {
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{},
	}

	return &HTTP{
		Client: client,
	}
}

// String returns a name for the client.
func (c *HTTP) String() string {
	return "http-client"
}

// Connect opens a long-lived connection to the external service.
func (c *HTTP) Connect() error {
	return nil
}

// Disconnect closes the long-lived connection to the external service.
func (c *HTTP) Disconnect(ctx context.Context) error {
	return nil
}

// HealthCheck checks the health of connection to the external service.
func (c *HTTP) HealthCheck(ctx context.Context) error {
	return nil
}
