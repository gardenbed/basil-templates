package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gardenbed/basil/graceful"
	"github.com/gardenbed/basil/health"
	"github.com/gardenbed/basil/httpx"

	githubentity "grpc-service-horizontal/internal/entity/github"
)

// Gateway is the interface for calling an external service.
type Gateway interface {
	graceful.Client
	health.Checker
	GetUser(ctx context.Context, username string) (*githubentity.User, error)
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// gateway implements the Gateway interface.
type gateway struct {
	client httpClient
}

// NewGateway creates a new gateway.
func NewGateway() (Gateway, error) {
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{},
	}

	return &gateway{
		client: client,
	}, nil
}

// String returns a name for the gateway.
func (g *gateway) String() string {
	return "github-gateway"
}

// Connect opens a long-lived connection to the external service.
func (g *gateway) Connect() error {
	return nil
}

// Disconnect closes the long-lived connection to the external service.
func (g *gateway) Disconnect(ctx context.Context) error {
	return nil
}

// HealthCheck checks the health of connection to the external service.
func (g *gateway) HealthCheck(ctx context.Context) error {
	return nil
}

// GetUser retrieves a GitHub user by username.
func (g *gateway) GetUser(ctx context.Context, username string) (*githubentity.User, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	req.Header.Set("User-Agent", "command-app")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, httpx.NewClientError(resp)
	}

	user := new(githubentity.User)
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}
