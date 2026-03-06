package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/oguzx/devpulse/internal/service"
)

type ServiceHandler struct {
	serviceService *service.ServiceService
}

func NewServiceHandler(s *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{
		serviceService: s,
	}
}

func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input service.CreateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	createService, err := h.serviceService.Create(r.Context(), input)

	if err != nil {
		return
	}

	writeJSON(w, http.StatusOK, createService)
}

func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	services, err := h.serviceService.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list services")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": services,
	})
}
