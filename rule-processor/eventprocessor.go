package rule_processor

import (
	"context"
	"github.com/SMART2016/go-rule-engine/models"
)

type RuleProcessor interface {
	Evaluate(ctx context.Context, event models.BaseEvent[any]) (bool, error)
}

type EventProcessor struct {
	ruleProc RuleProcessor
}

/*
*
This is the main API for the consumers of the framework to initialize the event processor.
This is a Go function named NewProcessor that creates a new instance of an event processor. It takes a configuration
object cfg as input and returns an EventProcessor instance and an error.
Here's a step-by-step breakdown:

It initializes a rule repository using the provided configuration.
It generates a database connection string (DSN) using the configuration's database settings.
It initializes an event store using the generated DSN.
It creates a new EventProcessor instance with a GRuleProcessor instance, which is configured with the initialized rule
repository and event store.
The function returns the EventProcessor instance and an error if any of the initialization steps fail.
*/
func NewProcessor(cfg Config) (*EventProcessor, error) {
	//Initialize Rule Processor
	ruleProc, err := NewGRuleProcessor(cfg)
	if err != nil {

	}
	// Create a new EventProcessor instance with a GRuleProcessor instance
	return &EventProcessor{
		ruleProc: ruleProc,
	}, nil
}

/*
ProcessEvent evaluates the event against the configured rules.

Parameters:
  - ctx: context.Context - Used to pass a context to the rule engine.
  - event: models.BaseEvent[any] - The event to be evaluated.

Returns:
  - bool: Indicates whether the event should be handled.
  - error: Error if any occurs during evaluation.

The function evaluates the event using the rule engine associated with the EventProcessor.
If the event should be handled, it returns true and nil error.
If the event should not be handled, it returns false and nil error.
If an error occurs during evaluation, it returns false and the error.
*/
func (ep *EventProcessor) ProcessEvent(ctx context.Context, event models.BaseEvent[any]) (bool, error) {
	// Evaluate the event using the rule engine
	return ep.ruleProc.Evaluate(ctx, event)
}
