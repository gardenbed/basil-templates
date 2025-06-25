package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gardenbed/basil/httpx"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Service is used for calling an external service.
type Service struct {
	client httpClient
}

// NewService creates a new service.
func NewService() (*Service, error) {
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{},
	}

	return &Service{
		client: client,
	}, nil
}

// User is the model for a GitHub user.
type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetUser retrieves a GitHub user by username.
func (s *Service) GetUser(ctx context.Context, username string) (*User, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	req.Header.Set("User-Agent", "command-line-app")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, httpx.NewClientError(resp)
	}

	user := new(User)
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}
