package translate

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGateway(t *testing.T) {
	tests := []struct {
		name          string
		expectedError string
	}{
		{
			name:          "OK",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGateway()

			if tc.expectedError == "" {
				assert.NotNil(t, g)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, g)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_String(t *testing.T) {
	g := new(gateway)
	assert.Equal(t, "translate-gateway", g.String())
}

func TestGateway_Connect(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockHTTPClient
		expectedError string
	}{
		{
			name:          "OK",
			client:        &MockHTTPClient{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			err := g.Connect()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_Disconnect(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockHTTPClient
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "OK",
			client:        &MockHTTPClient{},
			ctx:           context.Background(),
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			err := g.Disconnect(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_HealthCheck(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockHTTPClient
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "OK",
			client:        &MockHTTPClient{},
			ctx:           context.Background(),
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			err := g.HealthCheck(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGateway_Translate(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockHTTPClient
		ctx           context.Context
		lang          string
		text          string
		expectedText  string
		expectedError string
	}{
		{
			name: "DoFails",
			client: &MockHTTPClient{
				DoMocks: []DoMock{
					{OutError: errors.New("http error")},
				},
			},
			ctx:           context.Background(),
			lang:          "fr",
			text:          "Hello",
			expectedText:  "",
			expectedError: "http error",
		},
		{
			name: "InvalidResponseBody",
			client: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{`),
							),
						},
					},
				},
			},
			ctx:           context.Background(),
			lang:          "fr",
			text:          "Hello",
			expectedText:  "",
			expectedError: "unexpected EOF",
		},
		{
			name: "Success",
			client: &MockHTTPClient{
				DoMocks: []DoMock{
					{
						OutResponse: &http.Response{
							StatusCode: 200,
							Body: io.NopCloser(
								strings.NewReader(`{ "translatedText": "Bonjour" }`),
							),
						},
					},
				},
			},
			ctx:           context.Background(),
			lang:          "fr",
			text:          "Hello",
			expectedText:  "Bonjour",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateway{
				client: tc.client,
			}

			text, err := g.Translate(tc.ctx, tc.lang, tc.text)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedText, text)
			} else {
				assert.Empty(t, text)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
