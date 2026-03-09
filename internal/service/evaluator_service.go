package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/oguzx/devpulse/internal/domain"
	"github.com/oguzx/devpulse/internal/repository"
)

type EvaluatorService struct {
	serviceRepo  *repository.ServiceRepository
	incidentRepo *repository.IncidentRepository
	logger       *slog.Logger
}

func NewEvaluatorService(
	serviceRepo *repository.ServiceRepository,
	incidentRepo *repository.IncidentRepository,
	logger *slog.Logger,
) *EvaluatorService {
	return &EvaluatorService{
		serviceRepo:  serviceRepo,
		incidentRepo: incidentRepo,
		logger:       logger,
	}
}

func (s *EvaluatorService) EvaluateOnce(ctx context.Context) error {
	services, err := s.serviceRepo.FindAllForEvaluation(ctx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	s.logger.Info("evaluation cycle started", "service_count", len(services), "now", now.Format(time.RFC3339))

	for _, svc := range services {
		if svc.LastHeartbeatAt == nil {
			s.logger.Info(
				"skipping service without heartbeat",
				"service_id", svc.ID,
				"service_slug", svc.Slug,
			)
			continue
		}

		deadline := svc.LastHeartbeatAt.Add(
			time.Duration(svc.ExpectedIntervalSeconds+svc.GraceSeconds) * time.Second,
		)

		s.logger.Info(
			"evaluating service",
			"service_id", svc.ID,
			"service_slug", svc.Slug,
			"status", svc.Status,
			"last_heartbeat_at", svc.LastHeartbeatAt.Format(time.RFC3339),
			"deadline", deadline.Format(time.RFC3339),
			"now", now.Format(time.RFC3339),
		)

		if now.After(deadline) {
			s.logger.Info(
				"deadline exceeded",
				"service_id", svc.ID,
				"service_slug", svc.Slug,
			)

			if svc.Status != domain.ServiceStatusDown {
				if err := s.serviceRepo.MarkDown(ctx, svc.ID); err != nil {
					return err
				}

				s.logger.Info(
					"service marked down",
					"service_id", svc.ID,
					"service_slug", svc.Slug,
				)
			}

			_, err := s.incidentRepo.FindOpenByServiceID(ctx, svc.ID)
			if err == nil {
				s.logger.Info(
					"open incident already exists",
					"service_id", svc.ID,
					"service_slug", svc.Slug,
				)
				continue
			}

			if err != repository.ErrIncidentNotFound {
				return err
			}

			reason := fmt.Sprintf(
				"heartbeat deadline exceeded at %s",
				deadline.Format(time.RFC3339),
			)

			incident, err := s.incidentRepo.CreateOpen(ctx, svc.ID, deadline, reason)
			if err != nil {
				return err
			}

			s.logger.Info(
				"incident opened",
				"incident_id", incident.ID,
				"service_id", svc.ID,
				"service_slug", svc.Slug,
			)
		}
	}

	return nil
}

func (s *EvaluatorService) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.logger.Info("evaluator started", "interval", interval.String())
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("evaluator stopped")
			return
		case <-ticker.C:
			evalCtx, cancel := context.WithTimeout(ctx, time.Minute*10)
			err := s.EvaluateOnce(evalCtx)
			cancel()
			if err != nil {
				s.logger.Error("evaluator run failed", "error", err)
			}
		}
	}
}
