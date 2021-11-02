package handler

import (
	"context"

	"http-service-horizontal/internal/entity"
)

type GreetMock struct {
	InContext   context.Context
	InRequest   *entity.GreetRequest
	OutResponse *entity.GreetResponse
	OutError    error
}

// MockGreetingController is a mock implementation for controller.GreetingController.
type MockGreetingController struct {
	GreetIndex int
	GreetMocks []GreetMock
}

func (m *MockGreetingController) Greet(ctx context.Context, request *entity.GreetRequest) (*entity.GreetResponse, error) {
	i := m.GreetIndex
	m.GreetIndex++
	m.GreetMocks[i].InContext = ctx
	m.GreetMocks[i].InRequest = request
	return m.GreetMocks[i].OutResponse, m.GreetMocks[i].OutError
}
