package greeting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gardenbed/basil/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
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
	Name string `json:"name"`
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

	greeting, err := s.getGreeting(r.Context())
	if err != nil {
		httpx.Error(w, err, http.StatusInternalServerError)
		return
	}

	greeting = fmt.Sprintf("%s, %s!", greeting, req.Name)
	resp := &GreetResponse{
		Greeting: greeting,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Service) getGreeting(ctx context.Context) (string, error) {
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
