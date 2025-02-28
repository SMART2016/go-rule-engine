package models

import "time"

type Event[T any] struct {
	TenantID string `json:"tenant_id"`
	Type     string `json:"type"`
	//TODO: Handle dynamic payload schema
	Payload      T         `json:"payload"`
	ShouldHandle bool      `json:"should_handle"`
	EventSHA     string    `json:"event_sha"`
	OccuredAt    time.Time `json:"ocured_at"`
}
