package handler

import (
	"context"

	"grpc-service-horizontal/internal/controller/greeting"
	"grpc-service-horizontal/internal/idl/greetingpb"
	"grpc-service-horizontal/internal/mapper"
)

// GreetingHandler is an alias for the gRPC server interface.
type GreetingHandler = greetingpb.GreetingServiceServer

// greetingHandler implements GreetingHandler (greetingpb.GreetingServiceServer) interface.
type greetingHandler struct {
	greetingController greeting.Controller
}

// NewGreetingHandler creates a new instance of GreetingHandler.
func NewGreetingHandler(greetingController greeting.Controller) (GreetingHandler, error) {
	return &greetingHandler{
		greetingController: greetingController,
	}, nil
}

// Greet is the handler for GreetingService::Greet endpoint.
func (h *greetingHandler) Greet(ctx context.Context, req *greetingpb.GreetRequest) (*greetingpb.GreetResponse, error) {
	domainReq, err := mapper.GreetRequestIDLToDomain(req)
	if err != nil {
		return nil, err
	}

	domainResp, err := h.greetingController.Greet(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	resp, err := mapper.GreetResponseDomainToIDL(domainResp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
