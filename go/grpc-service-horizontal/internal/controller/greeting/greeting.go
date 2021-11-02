package greeting

import (
	"context"
	"fmt"

	"grpc-service-horizontal/internal/entity"
	"grpc-service-horizontal/internal/gateway/translate"
	"grpc-service-horizontal/internal/repository/greetingcache"
)

const lang = "fr"

// Controller is the interface for greeting business logic.
type Controller interface {
	Greet(context.Context, *entity.GreetRequest) (*entity.GreetResponse, error)
}

// controller implements the Controller interface.
type controller struct {
	translateGateway        translate.Gateway
	greetingcacheRepository greetingcache.Repository
}

// NewController creates a new controller.
func NewController(translateGateway translate.Gateway, greetingcacheRepository greetingcache.Repository) (Controller, error) {
	return &controller{
		translateGateway:        translateGateway,
		greetingcacheRepository: greetingcacheRepository,
	}, nil
}

// Greet creates and returns a greeting for a given name!
func (c *controller) Greet(ctx context.Context, req *entity.GreetRequest) (*entity.GreetResponse, error) {
	greeting, err := c.getGreeting(ctx)
	if err != nil {
		return nil, err
	}

	greeting = fmt.Sprintf("%s, %s!", greeting, req.Name)
	resp := &entity.GreetResponse{
		Greeting: greeting,
	}

	return resp, nil
}

func (c *controller) getGreeting(ctx context.Context) (string, error) {
	greeting, err := c.greetingcacheRepository.Lookup(ctx, lang)
	if err == nil && greeting != "" {
		return greeting, nil
	}

	greeting, err = c.translateGateway.Translate(ctx, lang, "Hello")
	if err != nil {
		return "", err
	}

	if greeting != "" {
		_ = c.greetingcacheRepository.Store(ctx, lang, "Hello")
		return greeting, nil
	}

	return "Hello", nil
}
