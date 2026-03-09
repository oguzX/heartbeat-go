package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/oguzx/devpulse/internal/domain"
	"github.com/oguzx/devpulse/internal/repository"
)

type HeartbeatService struct {
	serviceRepo   *repository.ServiceRepository
	heartbeatRepo *repository.HeartbeatRepository
	incidentRepo  *repository.IncidentRepository
}

func NewHeartbeatService(
	serviceRepo *repository.ServiceRepository,
	heartbeatRepo *repository.HeartbeatRepository,
	incidentRepo *repository.IncidentRepository,
) *HeartbeatService {
	return &HeartbeatService{
		serviceRepo:   serviceRepo,
		heartbeatRepo: heartbeatRepo,
		incidentRepo:  incidentRepo,
	}
}

type IngestHeartbeatInput struct {
	ServiceKey string          `json:"service_key"`
	Meta       json.RawMessage `json:"meta"`
}

type IngestHeartbeatResult struct {
	Service   *domain.Service   `json:"service"`
	Heartbeat *domain.Heartbeat `json:"heartbeat"`
}

func (s *HeartbeatService) Ingest(
	ctx context.Context,
	input IngestHeartbeatInput,
	r *http.Request,
) (*IngestHeartbeatResult, error) {
	serviceKey := strings.TrimSpace(input.ServiceKey)
	if serviceKey == "" {
		return nil, fmt.Errorf("service_key is required")
	}

	service, err := s.serviceRepo.FindByAPIKey(ctx, serviceKey)
	if err != nil {
		if err == repository.ErrServiceNotFound {
			return nil, fmt.Errorf("invalid service_key")
		}
		return nil, err
	}

	meta := input.Meta
	if len(meta) == 0 {
		meta = json.RawMessage(`{}`)
	}

	sourceIP := extractClientIP(r)

	heartbeat, err := s.heartbeatRepo.Create(ctx, service.ID, sourceIP, meta)
	if err != nil {
		return nil, err
	}

	if err := s.serviceRepo.MarkHealthy(ctx, service.ID); err != nil {
		return nil, err
	}

	if err := s.incidentRepo.ResolveOpenByServiceID(ctx, service.ID); err != nil {
		return nil, err
	}

	service, err = s.serviceRepo.FindByAPIKey(ctx, serviceKey)
	if err != nil {
		return nil, err
	}

	return &IngestHeartbeatResult{
		Service:   service,
		Heartbeat: heartbeat,
	}, nil
}

func extractClientIP(r *http.Request) *string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		value := strings.TrimSpace(parts[0])
		if value != "" {
			return &value
		}
	}

	realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if realIP != "" {
		return &realIP
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return &host
	}

	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	if remoteAddr != "" {
		return &remoteAddr
	}

	return nil
}
