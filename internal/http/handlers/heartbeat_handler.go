package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/oguzx/devpulse/internal/service"
)

type HeartbeatHandler struct {
	heartbeatService *service.HeartbeatService
}

func NewHeartbeatHandler(service *service.HeartbeatService) *HeartbeatHandler {
	return &HeartbeatHandler{
		heartbeatService: service,
	}
}

func (h *HeartbeatHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var input service.IngestHeartbeatInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.heartbeatService.Ingest(r.Context(), input, r)

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
