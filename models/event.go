package models

import (
	"errors"
	"time"
)

/*
Event represents a generic event structure with a dynamic payload.
It uses a type parameter T to allow for various payload schemas.

Fields:
- TenantID: Identifies the tenant associated with the event.
- Type: Specifies the type of the event.
- Payload: Contains the event-specific data, which can be of any type T.
- ShouldHandle: Indicates whether the event should be processed.
- EventSHA: A hash of the event data for deduplication purposes.
- OccuredAt: The timestamp when the event occurred.
*/
type Event[T any] struct {
	TenantID     string    `json:"tenant_id"` // cannot be empty, unmarshaling will fail if empty
	Type         string    `json:"type"`
	Payload      T         `json:"payload,omitempty"`
	ShouldHandle bool      `json:"should_handle,omitempty"`
	EventSHA     string    `json:"event_sha,omitempty"`
	OccuredAt    time.Time `json:"occured_at"`
}

// Validate checks required fields
func (e *Event[T]) Validate() error {
	if e.TenantID == "" {
		return errors.New("tenant_id cannot be empty")
	}
	if e.Type == "" {
		return errors.New("type cannot be empty")
	}
	return nil
}
