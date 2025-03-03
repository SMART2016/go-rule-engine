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
	//Initialize basic Configs
	config, err := ruleprocessor.NewFrameworkConfig(opts...)
	if err != nil {
		log.Fatalf("Error initializing framework config: %v", err)
	}

	// Initialize Rule Processor
	processor, err := ruleprocessor.NewGRuleProcessor(config)
	if err != nil {
		log.Fatalf("Error initializing rule processor: %v", err)
	}

	// Initialize event registry
	registry := models.GetEventRegistry()
	// Step 1: Register event types (only once at startup)
	registry.RegisterEventType("disk_space", func() models.Evaluable {
		return &examples.DiskUsageEvent{} // ✅ Correct type
	})

	// Example JSON input for a disk usage event
	//TODO: occurred_at needs to be set correctly in the event or use default event persistence time
	rawJSON := []byte(`{
		"payload": {
			"usage_percentage": 80,
			"instance_id": "abcd",
			"disk_size_in_bytes": 2048
		},
		"tenant_id": "tenant_1",
		"type": "disk_space"
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
