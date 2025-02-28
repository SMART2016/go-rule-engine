// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package store

import (
	"database/sql"

	"github.com/sqlc-dev/pqtype"
)

type ProcessedEvent struct {
	ID           int64                 `json:"id"`
	TenantID     string                `json:"tenant_id"`
	EventType    string                `json:"event_type"`
	EventSha     string                `json:"event_sha"`
	EventDetails pqtype.NullRawMessage `json:"event_details"`
	CreatedAt    sql.NullTime          `json:"created_at"`
}
