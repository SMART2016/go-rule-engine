package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// EventRegistry stores event constructors dynamically.
type EventRegistry struct {
	eventConstructors map[string]func() Evaluable
}

// Global registry instance.
var registry = &EventRegistry{
	eventConstructors: make(map[string]func() Evaluable),
}

func GetEventRegistry() *EventRegistry {
	return registry
}

// RegisterEventType registers an event type with a constructor function.
func (er *EventRegistry) GetRegistry() map[string]func() Evaluable {
	return registry.eventConstructors
}

// RegisterEventType registers an event type with a constructor function.
func (er *EventRegistry) RegisterEventType(eventType string, constructor func() Evaluable) {
	er.eventConstructors[eventType] = constructor
}

/*
ProcessEvent detects the event type and evaluates it.
ProcessEvent detects the event type from raw JSON data, constructs the corresponding event,
and evaluates it using the provided rule processor.

Parameters:
  - ctx: context.Context - A context to manage cancellation and deadlines.
  - processor: RuleProcessor - An interface for processing rules associated with the event.
  - rawJSON: []byte - The raw JSON data representing the event.

Returns:
  - bool - Indicates whether the event was handled successfully.
  - error - Contains any error encountered during processing or evaluation of the event.

NOTE: the consumer needs to handle the error and make sure the event that caused error while processing is
either logged properly or pushed into a dead letter queue.
*/
func (er *EventRegistry) ProcessEvent(ctx context.Context, processor RuleProcessor, rawJSON []byte) (bool, error) {
	// Step 1: Decode the event to extract the type field.
	var temp map[string]interface{}
	err := json.Unmarshal(rawJSON, &temp)
	if err != nil {
		return false, errors.New("failed to parse event JSON")
	}

	// Step 2: Extract the event type.
	eventType, ok := temp["type"].(string)
	if !ok {
		return false, errors.New("missing or invalid event type")
	}

	// Step 3: Look up the registered event constructor.
	constructor, found := registry.eventConstructors[eventType]
	if !found {
		return false, fmt.Errorf("event type '%s' not registered", eventType)
	}

	// Step 4: Create a new event instance using the constructor.
	eventInstance := constructor()

	// Step 5: Unmarshal JSON into the specific event struct.
	err = json.Unmarshal(rawJSON, eventInstance)
	if err != nil {
		return false, fmt.Errorf("failed to parse event payload: %v", err)
	}
	// Step 6: Evaluate the event.
	return eventInstance.Evaluate(ctx, processor)
}
