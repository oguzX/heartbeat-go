package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oguzx/devpulse/internal/domain"
)

type HeartBeatRepository struct {
	db *pgxpool.Pool
}

func NewHeartbeatRepository(db *pgxpool.Pool) *HeartBeatRepository {
	return &HeartBeatRepository{db: db}
}

func (r *HeartBeatRepository) Create(
	ctx context.Context,
	serviceID int64,
	sourceIP *string,
	meta json.RawMessage) (*domain.Heartbeat, error) {
	query := `
		INSERT INTO heartbeats (
			service_id,
			source_ip,
			meta_json
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			service_id,
			received_at,
			source_ip,
			meta_json
	`

	var heartbeat domain.Heartbeat
	err := r.db.QueryRow(ctx, query, serviceID, sourceIP, meta).Scan(
		&heartbeat.ID,
		&heartbeat.ServiceID,
		&heartbeat.ReceivedAt,
		&heartbeat.SourceIP,
		&heartbeat.MetaJSON,
	)

	if err != nil {
		return nil, err
	}

	return &heartbeat, nil
}
