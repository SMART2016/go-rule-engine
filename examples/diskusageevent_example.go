package examples

import (
	"context"
	"errors"
	"fmt"
	"github.com/SMART2016/go-rule-engine/models"
)

type RuleProcessor interface {
	Evaluate(ctx context.Context, event models.BaseEvent[any]) (bool, error)
}
type DiskUsagePayload struct {
	UsagePercentage int    `json:"usage_percentage"`
	InstanceID      string `json:"instance_id"`
	DiskSizeInBytes int64  `json:"disk_size_in_bytes"`
}

type DiskUsageEvent struct {
	models.BaseEvent[DiskUsagePayload]
}

/*
Validate checks required fields

Returns an error if any of the required fields are invalid or empty.
*/
func (e *DiskUsageEvent) Validate() error {
	// This should be called by all custom events to handle basic field validations
	e.BaseEvent.Validate()

	if e.Payload.UsagePercentage < 0 {
		return errors.New("usage_percentage should be greater than or equal to 0")
	}
	if e.Payload.InstanceID == "" {
		return errors.New("instance_id cannot be empty")
	}

	if e.Payload.DiskSizeInBytes <= 0 {
		return errors.New("disk_size_in_bytes should be greater than 0")
	}

	return nil
}

/*
DeduplicationKeyValues directly returns values instead of field names.
DeduplicationKeyValues generates a pipe seperated string that can be used as a deduplication key.

It uses the "usage_percentage" and "instance_id" fields in this example event from the event payload
to generate a string that is unique to the event. This string is then used
as the deduplication key to prevent duplicate events from being processed.
*/
func (e *DiskUsageEvent) DeduplicationKeyValues() string {
	return fmt.Sprintf("%d|%s", e.Payload.UsagePercentage, e.Payload.InstanceID)
}

// Implement Evaluate function
func (e *DiskUsageEvent) Evaluate(ctx context.Context, processor RuleProcessor) (bool, error) {
	// Generate SHA256 before passing event
	var err error
	e.EventSHA, err = e.GenerateSHA256(e.DeduplicationKeyValues())
	if err != nil {
		return false, err
	}
	// Pass self as BaseEvent[any]
	return processor.Evaluate(ctx, models.BaseEvent[any]{
		TenantID:     e.TenantID,
		Type:         e.Type,
		Payload:      e.Payload, // Convert to `any`
		ShouldHandle: e.ShouldHandle,
		EventSHA:     e.EventSHA,
	})
}
