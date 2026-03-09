package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oguzx/devpulse/internal/domain"
)

var ErrIncidentNotFound = errors.New("incident not found")

type IncidentRepository struct {
	db *pgxpool.Pool
}

func NewIncidentRepository(db *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) FindOpenByServiceID(ctx context.Context, serviceID int64) (*domain.Incident, error) {
	query := `
	SELECT
			id,
			service_id,
			started_at,
			resolved_at,
			status,
			reason,
			created_at,
			updated_at
		FROM incidents
		WHERE service_id = $1
		  AND status = $2
		ORDER BY started_at DESC
		LIMIT 1
		`

	var incident domain.Incident
	err := r.db.QueryRow(ctx, query, serviceID, domain.IncidentStatusOpen).Scan(
		&incident.ID,
		&incident.ServiceID,
		&incident.StartedAt,
		&incident.ResolvedAt,
		&incident.Status,
		&incident.Reason,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIncidentNotFound
		}

		return nil, err
	}

	return &incident, nil
}

func (r *IncidentRepository) CreateOpen(ctx context.Context, serviceID int64, started_at any, reason string) (*domain.Incident, error) {
	query := `INSERT INTO incidents (
			service_id,
			started_at,
			status,
			reason
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			service_id,
			started_at,
			resolved_at,
			status,
			reason,
			created_at,
			updated_at`

	var incident domain.Incident
	err := r.db.QueryRow(ctx, query, serviceID, started_at, domain.IncidentStatusOpen, reason).Scan(
		&incident.ID,
		&incident.ServiceID,
		&incident.StartedAt,
		&incident.ResolvedAt,
		&incident.Status,
		&incident.Reason,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &incident, nil
}

func (r *IncidentRepository) ResolveOpenByServiceID(ctx context.Context, serviceID int64) error {
	query := `
		UPDATE incidents
		SET
			status = $2,
			resolved_at = NOW(),
			updated_at = NOW()
		WHERE service_id = $1
		  AND status = $3
	`
	_, err := r.db.Exec(
		ctx,
		query,
		serviceID,
		domain.IncidentStatusResolved,
		domain.IncidentStatusOpen,
	)

	return err
}

func (r *IncidentRepository) List(ctx context.Context) ([]domain.Incident, error) {
	query := `
		SELECT
			id,
			service_id,
			started_at,
			resolved_at,
			status,
			reason,
			created_at,
			updated_at
		FROM incidents
		ORDER BY started_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incidents := make([]domain.Incident, 0)
	for rows.Next() {
		var incident domain.Incident
		if err := rows.Scan(

			&incident.ID,
			&incident.ServiceID,
			&incident.StartedAt,
			&incident.ResolvedAt,
			&incident.Status,
			&incident.Reason,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		); err != nil {
			return nil, err
		}

		incidents = append(incidents, incident)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}
