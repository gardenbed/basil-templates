package greeting

import (
	"context"
	"fmt"

	"http-service-horizontal/internal/entity"
	"http-service-horizontal/internal/gateway/github"
	"http-service-horizontal/internal/repository/usercache"
)

// Controller is the interface for greeting business logic.
type Controller interface {
	Greet(context.Context, *entity.GreetRequest) (*entity.GreetResponse, error)
}

// controller implements the Controller interface.
type controller struct {
	githubGateway       github.Gateway
	usercacheRepository usercache.Repository
}

// NewController creates a new controller.
func NewController(githubGateway github.Gateway, usercacheRepository usercache.Repository) (Controller, error) {
	return &controller{
		githubGateway:       githubGateway,
		usercacheRepository: usercacheRepository,
	}, nil
}

// Greet creates a greeting for a given GitHub user!
func (c *controller) Greet(ctx context.Context, req *entity.GreetRequest) (*entity.GreetResponse, error) {
	name, err := c.getName(ctx, req.GithubUsername)
	if err != nil {
		return nil, err
	}

	greeting := fmt.Sprintf("Hello, %s!", name)
	resp := &entity.GreetResponse{
		Greeting: greeting,
	}

	return resp, nil
}

func (c *controller) getName(ctx context.Context, username string) (string, error) {
	name, err := c.usercacheRepository.Lookup(ctx, username)
	if err == nil && name != "" {
		return name, nil
	}

	user, err := c.githubGateway.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	_ = c.usercacheRepository.Store(ctx, username, user.Name)

	return user.Name, nil
}
