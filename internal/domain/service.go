package domain

import "time"

type ServiceStatus string

const (
	ServiceStatusUnknown ServiceStatus = "unknown"
	ServiceStatusHealthy ServiceStatus = "healthy"
	ServiceStatusDown    ServiceStatus = "down"
)

type Service struct {
	ID                      int64         `json:"id"`
	Name                    string        `json:"name"`
	Slug                    string        `json:"slug"`
	APIKey                  string        `json:"api_key,omitempty"`
	ExpectedIntervalSeconds int32         `json:"expected_interval_seconds"`
	GraceSeconds            int32         `json:"grace_seconds"`
	Status                  ServiceStatus `json:"status"`
	LastHeartbeatAt         *time.Time    `json:"last_heartbeat_at,omitempty"`
	CreatedAt               time.Time     `json:"created_at"`
	UpdatedAt               time.Time     `json:"updated_at"`
}
