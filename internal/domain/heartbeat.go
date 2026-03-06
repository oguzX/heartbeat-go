package domain

import (
	"encoding/json"
	"time"
)

type Heartbeat struct {
	ID         int64           `json:"id"`
	ServiceID  int64           `json:"service_id"`
	ReceivedAt time.Time       `json:"received_at"`
	SourceIP   *string         `json:"source_ip,omitempty"`
	MetaJSON   json.RawMessage `json:"meta_json"`
}
