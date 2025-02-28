package rule_processor

import (
	"context"
	"github.com/SMART2016/go-rule-engine/models"
)

type RuleEvaluator interface {
	Evaluate(ctx context.Context, event interface{}) (bool, error)
}

type EventProcessor[T any] struct {
	ruleEngine RuleEvaluator
}

func NewEventProcessor[T any](ruleEngine RuleEvaluator) *EventProcessor[T] {
	return &EventProcessor[T]{
		ruleEngine: ruleEngine,
	}
}

func (ep *EventProcessor[T]) ProcessEvent(ctx context.Context, event models.Event[T]) (bool, error) {
	// Evaluate the event using the rule engine
	return ep.ruleEngine.Evaluate(ctx, event)
}
