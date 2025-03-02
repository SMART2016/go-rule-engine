package examples

import (
	"context"
	"errors"
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
	if e.TenantID == "" {
		return errors.New("tenant_id cannot be empty")
	}
	if e.Type == "" {
		return errors.New("type cannot be empty")
	}

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

// Implement Evaluate function
func (e *DiskUsageEvent) Evaluate(ctx context.Context, processor RuleProcessor) (bool, error) {
	// Pass self as BaseEvent[any]
	return processor.Evaluate(ctx, models.BaseEvent[any]{
		TenantID:     e.TenantID,
		Type:         e.Type,
		Payload:      e.Payload, // Convert to `any`
		ShouldHandle: e.ShouldHandle,
		EventSHA:     e.EventSHA,
	})
}
