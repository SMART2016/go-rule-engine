package main

import (
	"context"
	"fmt"
	"github.com/SMART2016/go-rule-engine/examples"
	"github.com/SMART2016/go-rule-engine/models"
	ruleprocessor "github.com/SMART2016/go-rule-engine/rule-processor"
	"log"
	"time"
)

func main() {
	//Initialize the config instance
	opts := []ruleprocessor.FrameworkConfigOption{
		ruleprocessor.WithDBConfigPath("configs/db_config.json"),
		ruleprocessor.WithRuleRepoPath("configs/rules.json"),
		ruleprocessor.WithCleanupInterval(1 * time.Hour),
	}
	config, err := ruleprocessor.NewFrameworkConfig(opts...)
	if err != nil {
		log.Fatalf("Error initializing framework config: %v", err)
	}

	// Initialize dependencies
	processor, err := ruleprocessor.NewGRuleProcessor(config)
	if err != nil {
		log.Fatalf("Error initializing rule processor: %v", err)
	}

	// Initialize event registry
	registry := ruleprocessor.GetEventRegistry()
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
