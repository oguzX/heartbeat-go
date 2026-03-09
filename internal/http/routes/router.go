package routes

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oguzx/devpulse/internal/http/handlers"
	"github.com/oguzx/devpulse/internal/repository"
	appservice "github.com/oguzx/devpulse/internal/service"
)

type AppDependencies struct {
	Router    http.Handler
	Evaluator *appservice.EvaluatorService
}

func NewApp(db *pgxpool.Pool, logger *slog.Logger) *AppDependencies {
	r := chi.NewRouter()

	healthHandler := handlers.NewHealthHandler(db)

	serviceRepo := repository.NewServiceRepository(db)
	heartbeatRepo := repository.NewHeartbeatRepository(db)
	incidentRepo := repository.NewIncidentRepository(db)

	serviceService := appservice.NewServiceService(serviceRepo)
	heartbeatService := appservice.NewHeartbeatService(serviceRepo, heartbeatRepo, incidentRepo)
	evaluatorService := appservice.NewEvaluatorService(serviceRepo, incidentRepo, logger)

	serviceHandler := handlers.NewServiceHandler(serviceService)
	heartbeatHandler := handlers.NewHeartbeatHandler(heartbeatService)
	incidentHandler := handlers.NewIncidentHandler(incidentRepo)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("DevPulse is running"))
	})

	r.Get("/health", healthHandler.ServerHTTP)

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/services", func(r chi.Router) {
			r.Post("/", serviceHandler.Create)
			r.Get("/", serviceHandler.List)
		})

		r.Post("/heartbeats", heartbeatHandler.Ingest)
		r.Get("/incidents", incidentHandler.List)
	})

	return &AppDependencies{
		Router:    r,
		Evaluator: evaluatorService,
	}
}
