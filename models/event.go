package models

import (
	"context"
	"errors"
	"time"
)

type RuleProcessor interface {
	// Evaluate takes a context and a BaseEvent and returns a boolean indicating
	// whether the event was handled and an error if there was a problem.
	Evaluate(ctx context.Context, event BaseEvent[any]) (bool, error)
}

// Evaluable ensures all events implement `Evaluate`
type Evaluable interface {
	Evaluate(ctx context.Context, processor RuleProcessor) (bool, error)
}

/*
BaseEvent represents a generic event structure with a dynamic payload.
It uses a type parameter T to allow for various payload schemas.

Fields:
- TenantID: Identifies the tenant associated with the event.
- Type: Specifies the type of the event.
- Payload: Contains the event-specific data, which can be of any type T.
- ShouldHandle: Indicates whether the event should be processed.
- EventSHA: A hash of the event data for deduplication purposes.
- OccuredAt: The timestamp when the event occurred.
*/
type BaseEvent[T any] struct {
	TenantID     string    `json:"tenant_id"` // cannot be empty, unmarshaling will fail if empty
	Type         string    `json:"type"`
	Payload      T         `json:"payload,omitempty"`
	ShouldHandle bool      `json:"should_handle,omitempty"`
	EventSHA     string    `json:"event_sha,omitempty"`
	OccuredAt    time.Time `json:"occured_at"`
}

func (e *BaseEvent[T]) Evaluate(ctx context.Context, processor RuleProcessor) (bool, error) {
	//TODO implement me
	return processor.Evaluate(ctx, BaseEvent[any]{
		TenantID:     e.TenantID,
		Type:         e.Type,
		Payload:      e.Payload, // Convert to `any`
		ShouldHandle: e.ShouldHandle,
		EventSHA:     e.EventSHA,
	})
}

// Validate checks required fields
func (e *BaseEvent[T]) Validate() error {
	if e.TenantID == "" {
		return errors.New("tenant_id cannot be empty")
	}
	if e.Type == "" {
		return errors.New("type cannot be empty")
	}
	return nil
}
