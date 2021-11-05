package greeting

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gardenbed/basil/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type (
	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	redisClient interface {
		Get(context.Context, string) *redis.StringCmd
		Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	}
)

// Service implements the HTTP handlers for Greeting APIs.
type Service struct {
	httpClient  httpClient
	redisClient redisClient
}

// NewService creates a new service.
func NewService(httpClient httpClient, redisClient redisClient) (*Service, error) {
	return &Service{
		httpClient:  httpClient,
		redisClient: redisClient,
	}, nil
}

// RegisterRoutes registers the HTTP routes for greeting service.
// Middleware are applied from left to right (the first middleware is the most inner and the last middleware is the most outter).
func (s *Service) RegisterRoutes(router *mux.Router, middleware ...httpx.Middleware) {
	greetHandler := s.Greet

	for _, mid := range middleware {
		greetHandler = mid.Wrap(greetHandler)
	}

	router.Name("Greet").Methods("POST").Path("/v1/greet").HandlerFunc(greetHandler)
}

// GreetRequest is the HTTP (wire/transport protocol) model for a Greet request.
type GreetRequest struct {
	GithubUsername string `json:"githubUsername"`
}

// GreetResponse is the HTTP (wire/transport protocol) model for a Greet response.
type GreetResponse struct {
	Greeting string `json:"greeting"`
}

// Greet is the handler for the GreetingService::Greet endpoint.
func (s *Service) Greet(w http.ResponseWriter, r *http.Request) {
	req := new(GreetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	name, err := s.getName(r.Context(), req.GithubUsername)
	if err != nil {
		httpx.Error(w, err, http.StatusInternalServerError)
		return
	}

	greeting := fmt.Sprintf("Hello, %s!", name)
	resp := &GreetResponse{
		Greeting: greeting,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Service) getName(ctx context.Context, username string) (string, error) {
	name, err := s.redisClient.Get(ctx, username).Result()
	if err == nil && name != "" {
		return name, nil
	}

	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	req.Header.Set("User-Agent", "command-app")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", httpx.NewClientError(resp)
	}

	user := struct {
		Name string `json:"name"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	return user.Name, nil
}
