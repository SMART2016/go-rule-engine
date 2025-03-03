package models

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

	switch payload := any(e.Payload).(type) {
	case string:
		if !json.Valid([]byte(payload)) {
			return errors.New("payload is a string but not valid JSON")
		}
	case struct{}: // Ensures it's a struct
		// Struct is valid, do nothing
	default:
		return fmt.Errorf("Invalid payload type: expected struct or JSON string, got %T", e.Payload)
	}
	return nil
}

/*
GenerateSHA256 computes a SHA256 hash from deduplication key values.

This function takes a string containing deduplication key values,
computes its SHA256 hash, and returns the hash encoded as a hexadecimal string.
The deduplication keys should be a unique representation of the event data
to ensure proper deduplication.
*/
func (e *BaseEvent[T]) GenerateSHA256(dedupKeys string) (string, error) {
	if dedupKeys == "" {
		return "", errors.New("deduplication key values cannot be empty")
	}
	hash := sha256.Sum256([]byte(dedupKeys))
	return hex.EncodeToString(hash[:]), nil
}

// ToJSON converts the payload to a JSON string for storage in PostgreSQL.
func (e *BaseEvent[T]) ToJSON(any) (string, error) {
	// Check if payload is already a JSON string
	if jsonStr, ok := any(e.Payload).(string); ok {
		if json.Valid([]byte(jsonStr)) {
			return jsonStr, nil // Return as-is if it's valid JSON
		}
	}

	// Convert payload to JSON
	payloadJSON, err := json.Marshal(e.Payload)
	if err != nil {
		return "", fmt.Errorf("failed to convert payload to JSON: %v", err)
	}

	return string(payloadJSON), nil
}
