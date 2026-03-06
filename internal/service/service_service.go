package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/oguzx/devpulse/internal/domain"
	"github.com/oguzx/devpulse/internal/repository"
)

type ServiceService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(repo *repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

type CreateServiceInput struct {
	Name                    string `json:"name"`
	Slug                    string `json:"slug"`
	ExpectedIntervalSeconds int32  `json:"expected_interval_seconds"`
	GraceSeconds            int32  `json:"grace_seconds"`
}

func (s *ServiceService) Create(ctx context.Context, input CreateServiceInput) (*domain.Service, error) {
	name := strings.TrimSpace(input.Name)
	slug := normalizeSlug(input.Slug)

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if slug == "" {
		return nil, fmt.Errorf("slug is required")
	}

	if !isValidSlug(slug) {
		return nil, fmt.Errorf("slug must contain only lowercase letters, numbers, and hyphens")
	}

	if input.ExpectedIntervalSeconds <= 0 {
		input.ExpectedIntervalSeconds = 60
	}

	if input.GraceSeconds < 0 {
		input.GraceSeconds = 30
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate api key: %w", err)
	}

	service, err := s.repo.Create(
		ctx,
		name,
		slug,
		apiKey,
		input.ExpectedIntervalSeconds,
		input.GraceSeconds,
	)

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) List(ctx context.Context) ([]domain.Service, error) {
	return s.repo.List(ctx)
}

func normalizeSlug(slug string) string {
	slug = strings.TrimSpace(strings.ToLower(slug))
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func isValidSlug(value string) bool {
	re := regexp.MustCompile(`^[a-z0-9-]+$`)
	return re.MatchString(value)
}

func generateAPIKey() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return "svc_" + hex.EncodeToString(buf), nil
}
