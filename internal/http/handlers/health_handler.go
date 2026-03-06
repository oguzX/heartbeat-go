package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{DB: db}
}

func (h *HealthHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	dbStatus := "up"

	if err := h.DB.Ping(r.Context()); err != nil {
		dbStatus = "down"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")

	response := healthResponse{
		Status:   "ok",
		Database: dbStatus,
	}

	_ = json.NewEncoder(w).Encode(response)
}
