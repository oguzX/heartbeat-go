package handlers

import (
	"net/http"

	"github.com/oguzx/devpulse/internal/repository"
)

type IncidentHandler struct {
	incidentRepo *repository.IncidentRepository
}

func NewIncidentHandler(repository *repository.IncidentRepository) *IncidentHandler {
	return &IncidentHandler{
		incidentRepo: repository,
	}
}

func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	incidents, err := h.incidentRepo.List(r.Context())

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list incidents")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": incidents,
	})
}
