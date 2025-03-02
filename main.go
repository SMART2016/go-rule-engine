package main

import (
	"context"
	"fmt"
	"github.com/SMART2016/go-rule-engine/examples"
	"github.com/SMART2016/go-rule-engine/models"
	ruleprocessor "github.com/SMART2016/go-rule-engine/rule-processor"
	"log"
)

func main() {
	//TODO: Initialize the config instance

	// Initialize dependencies
	processor, err := ruleprocessor.NewGRuleProcessor(nil)

	// Initialize event registry
	registry := ruleprocessor.getEventRegistry()
	// Step 1: Register event types (only once at startup)
	registry.RegisterEventType("disk_usage", func() models.Evaluable {
		return &models.BaseEvent[any]{
			TenantID:     "",
			Type:         "",
			Payload:      examples.DiskUsageEvent{},
			ShouldHandle: false,
			EventSHA:     "",
		}
	})

	// Example JSON input for a disk usage event
	rawJSON := []byte(`{
		"payload": {
			"usage_percentage": 80,
			"instance_id": "abcd",
			"disk_size_in_bytes": 2048
		},
		"tenant_id": "12345A",
		"type": "disk_usage"
	}`)

	// Step 2: Process the event dynamically
	ctx := context.Background()
	handled, err := registry.ProcessEvent(ctx, processor, rawJSON)
	if err != nil {
		log.Fatalf("Error processing event: %v", err)
	}
	if handled {
		fmt.Println("Event processed, action triggered.")
	} else {
		fmt.Println("Event processed, no action needed.")
	}
}
