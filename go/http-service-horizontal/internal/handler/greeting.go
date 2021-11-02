package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gardenbed/basil/httpx"

	"http-service-horizontal/internal/controller/greeting"
	"http-service-horizontal/internal/idl"
	"http-service-horizontal/internal/mapper"
)

// GreetingHandler is an alias for the HTTP service interface.
type GreetingHandler = idl.GreetingHandler

// greetingHandler implements GreetingHandler (idl.GreetingHandler) interface.
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
func (h *greetingHandler) Greet(w http.ResponseWriter, r *http.Request) {
	req := new(idl.GreetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	domainReq, err := mapper.GreetRequestIDLToDomain(req)
	if err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	domainResp, err := h.greetingController.Greet(r.Context(), domainReq)
	if err != nil {
		httpx.Error(w, err, http.StatusInternalServerError)
		return
	}

	resp, err := mapper.GreetResponseDomainToIDL(domainResp)
	if err != nil {
		httpx.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
