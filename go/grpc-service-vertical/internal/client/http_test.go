package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTP(t *testing.T) {
	c := NewHTTP()
	assert.NotNil(t, c)
}

func TestHTTP_String(t *testing.T) {
	c := &HTTP{}
	str := c.String()

	assert.Equal(t, "http-client", str)
}

func TestHTTP_Connect(t *testing.T) {
	c := &HTTP{}
	err := c.Connect()

	assert.NoError(t, err)
}

func TestHTTP_Disconnect(t *testing.T) {
	c := &HTTP{}
	ctx := context.Background()
	err := c.Disconnect(ctx)

	assert.NoError(t, err)
}

func TestGateway_HealthCheck(t *testing.T) {
	c := &HTTP{}
	ctx := context.Background()
	err := c.HealthCheck(ctx)

	assert.NoError(t, err)
}
