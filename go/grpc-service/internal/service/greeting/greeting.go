package greeting

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gardenbed/basil/httpx"
	"github.com/redis/go-redis/v9"

	"grpc-service/internal/idl/greetingpb"
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

// service implements the greetingpb.GreetingServiceServer interface.
type service struct {
	httpClient  httpClient
	redisClient redisClient
}

// NewService creates a new service.
func NewService(httpClient httpClient, redisClient redisClient) (greetingpb.GreetingServiceServer, error) {
	return &service{
		httpClient:  httpClient,
		redisClient: redisClient,
	}, nil
}

// Greet implements the GreetingService::Greet endpoint.
func (s *service) Greet(ctx context.Context, req *greetingpb.GreetRequest) (*greetingpb.GreetResponse, error) {
	name, err := s.getName(ctx, req.GithubUsername)
	if err != nil {
		return nil, err
	}

	greeting := fmt.Sprintf("Hello, %s!", name)
	resp := &greetingpb.GreetResponse{
		Greeting: greeting,
	}

	return resp, nil
}

func (s *service) getName(ctx context.Context, username string) (string, error) {
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

	defer func() {
		_ = resp.Body.Close()
	}()

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
