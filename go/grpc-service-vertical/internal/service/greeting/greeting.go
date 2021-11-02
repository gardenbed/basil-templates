package greeting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"

	"grpc-service-vertical/internal/idl/greetingpb"
)

const lang = "fr"

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
	greeting, err := s.getGreeting(ctx)
	if err != nil {
		return nil, err
	}

	greeting = fmt.Sprintf("%s, %s!", greeting, req.Name)
	resp := &greetingpb.GreetResponse{
		Greeting: greeting,
	}

	return resp, nil
}

func (s *service) getGreeting(ctx context.Context) (string, error) {
	greeting, err := s.redisClient.Get(ctx, lang).Result()
	if err == nil && greeting != "" {
		return greeting, nil
	}

	reqBody := struct {
		Query  string `json:"q"`
		Source string `json:"source"`
		Target string `json:"target"`
		Format string `json:"format"`
	}{
		Query:  "Hello",
		Source: "en",
		Target: lang,
		Format: "text",
	}

	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://libretranslate.com/translate", buf)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody := struct {
		Text string `json:"translatedText"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}

	greeting = respBody.Text

	if greeting != "" {
		_ = s.redisClient.Set(ctx, lang, "Hello", 0).Err()
		return greeting, nil
	}

	return "Hello", nil
}
