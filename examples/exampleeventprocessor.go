package examples

import (
	"context"
	"fmt"
	"github.com/SMART2016/go-rule-engine/examples/events"
	"github.com/SMART2016/go-rule-engine/models"
	ruleprocessor "github.com/SMART2016/go-rule-engine/rule-processor"
	"log"
	"time"
)

// ExampleRuleProcessor shows how to use the rule processor framework.
// ExampleRuleProcessor demonstrates how to use the rule processor framework.
//
// First, it initializes the framework config using the With* functions.
// Then, it initializes the rule processor using the NewGRuleProcessor function.
// Next, it registers an event type with the event registry using the RegisterEventType function.
// After that, it creates a JSON byte slice representing the event payload.
// Finally, it calls the ProcessEvent function to evaluate the event and trigger an action if needed.
//
// The example code is a self-contained demonstration of the rule processor framework.
// It is suitable for use in a Go test or example.
func ExampleRuleProcessor() {
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
		return &events.DiskUsageEvent{} // Correct type
	})

	// Example JSON input for a disk usage event
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
	// Output:
	// Event processed, action triggered.
}
