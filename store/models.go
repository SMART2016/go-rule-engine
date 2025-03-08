// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package store

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ProcessedEvent struct {
	ID                          int64
	TenantID                    string
	EventType                   string
	RuleID                      string
	EventSha                    string
	EventDetails                []byte
	OccurredAt                  pgtype.Timestamp
	ActualEventPersistentceTime pgtype.Timestamp
}
