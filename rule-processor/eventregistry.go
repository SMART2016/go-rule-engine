package rule_processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SMART2016/go-rule-engine/models"
)

// EventRegistry stores event constructors dynamically.
type EventRegistry struct {
	eventConstructors map[string]func() models.Evaluable
}

// Global registry instance.
var registry = &EventRegistry{
	eventConstructors: make(map[string]func() models.Evaluable),
}

func GetEventRegistry() *EventRegistry {
	return registry
}

// RegisterEventType registers an event type with a constructor function.
func (er *EventRegistry) GetRegistry() map[string]func() models.Evaluable {
	return registry.eventConstructors
}

// RegisterEventType registers an event type with a constructor function.
func (er *EventRegistry) RegisterEventType(eventType string, constructor func() models.Evaluable) {
	er.eventConstructors[eventType] = constructor
}

// ProcessEvent detects the event type and evaluates it.
func (er *EventRegistry) ProcessEvent(ctx context.Context, processor RuleProcessor, rawJSON []byte) (bool, error) {
	// Step 1: Decode event to get the type field
	var temp map[string]interface{}
	err := json.Unmarshal(rawJSON, &temp)
	if err != nil {
		return false, errors.New("failed to parse event JSON")
	}

	// Step 2: Extract event type
	eventType, ok := temp["type"].(string)
	if !ok {
		return false, errors.New("missing or invalid event type")
	}

	// Step 3: Look up registered event constructor
	constructor, found := registry.eventConstructors[eventType]
	if !found {
		return false, fmt.Errorf("event type '%s' not registered", eventType)
	}

	// Step 4: Create a new event instance using the constructor
	eventInstance := constructor()

	// Step 5: Unmarshal JSON into the specific event struct
	err = json.Unmarshal(rawJSON, eventInstance)
	if err != nil {
		return false, fmt.Errorf("failed to parse event payload: %v", err)
	}

	// Step 6: Evaluate the event
	return eventInstance.Evaluate(ctx, processor)
}
