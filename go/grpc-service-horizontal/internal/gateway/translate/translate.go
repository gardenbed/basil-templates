package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gardenbed/basil/graceful"
	"github.com/gardenbed/basil/health"
)

// Gateway is the interface for calling an external translation service.
type Gateway interface {
	graceful.Client
	health.Checker
	Translate(ctx context.Context, lang, text string) (string, error)
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
	return "translate-gateway"
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

// Translate calls an external API to translate a text into another language.
func (g *gateway) Translate(ctx context.Context, lang, text string) (string, error) {
	reqBody := struct {
		Query  string `json:"q"`
		Source string `json:"source"`
		Target string `json:"target"`
		Format string `json:"format"`
	}{
		Query:  text,
		Source: "en",
		Target: lang,
		Format: "text",
	}

	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://libretranslate.com/translate", buf)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
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

	return respBody.Text, nil
}
