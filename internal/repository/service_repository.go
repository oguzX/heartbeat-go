package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oguzx/devpulse/internal/domain"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(
	ctx context.Context,
	name string,
	slug string,
	apiKey string,
	expectedIntervalSeconds int32,
	graceSeconds int32,
) (*domain.Service, error) {
	query := `
		INSERT INTO services (
			name,
			slug,
			api_key,
			expected_interval_seconds,
			grace_seconds,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			name,
			slug,
			api_key,
			expected_interval_seconds,
			grace_seconds,
			status,
			last_heartbeat_at,
			created_at,
			updated_at
	`
	var service domain.Service
	err := r.db.QueryRow(
		ctx,
		query,
		name,
		slug,
		apiKey,
		expectedIntervalSeconds,
		graceSeconds,
		domain.ServiceStatusUnknown,
	).Scan(
		&service.ID,
		&service.Name,
		&service.Slug,
		&service.APIKey,
		&service.ExpectedIntervalSeconds,
		&service.GraceSeconds,
		&service.Status,
		&service.LastHeartbeatAt,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return nil, fmt.Errorf("service with same slug or api key already exists")
		}

		return nil, err
	}

	return &service, nil
}

func (r *ServiceRepository) List(ctx context.Context) ([]domain.Service, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			api_key,
			expected_interval_seconds,
			grace_seconds,
			status,
			last_heartbeat_at,
			created_at,
			updated_at
		FROM services
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := make([]domain.Service, 0)
	for rows.Next() {
		var service domain.Service
		if err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Slug,
			&service.APIKey,
			&service.ExpectedIntervalSeconds,
			&service.GraceSeconds,
			&service.Status,
			&service.LastHeartbeatAt,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *ServiceRepository) FindByAPIKey(ctx context.Context, apiKey string) (*domain.Service, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			api_key,
			expected_interval_seconds,
			grace_seconds,
			status,
			last_heartbeat_at,
			created_at,
			updated_at
		FROM services
		WHERE api_key = $1
		LIMIT 1
	`

	var service domain.Service
	err := r.db.QueryRow(ctx, query, apiKey).Scan(
		&service.ID,
		&service.Name,
		&service.Slug,
		&service.APIKey,
		&service.ExpectedIntervalSeconds,
		&service.GraceSeconds,
		&service.Status,
		&service.LastHeartbeatAt,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no rows") {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}

	return &service, nil
}

func (r *ServiceRepository) MarkHealthy(ctx context.Context, serviceID int64) error {
	query := `
		UPDATE services
		SET
			status = $2,
			last_heartbeat_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, serviceID, domain.ServiceStatusHealthy)
	return err
}
