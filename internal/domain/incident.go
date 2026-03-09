package domain

import "time"

type IncidentStatus string

const (
	IncidentStatusOpen     IncidentStatus = "open"
	IncidentStatusResolved IncidentStatus = "resolved"
)

type Incident struct {
	ID         int64          `json:"id"`
	ServiceID  int64          `json:"service_id"`
	StartedAt  time.Time      `json:"started_at"`
	ResolvedAt *time.Time     `json:"resolved_at,omitempty"`
	Status     IncidentStatus `json:"status"`
	Reason     string         `json:"reason"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
