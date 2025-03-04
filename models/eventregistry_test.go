package models

import (
	"context"
	"testing"
)

// MockEvaluable implements Evaluable interface for testing
type MockEvaluable struct {
	EvaluateFunc func(ctx context.Context, processor RuleProcessor) (bool, error)
}

func (m *MockEvaluable) Evaluate(ctx context.Context, processor RuleProcessor) (bool, error) {
	return m.EvaluateFunc(ctx, processor)
}

// MockRuleProcessor implements RuleProcessor interface for testing
type MockRuleProcessor struct{}

func (m MockRuleProcessor) Evaluate(ctx context.Context, event BaseEvent[any]) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func TestGetEventRegistry(t *testing.T) {
	reg := GetEventRegistry()
	if reg == nil {
		t.Error("GetEventRegistry returned nil")
	}
	if reg.eventConstructors == nil {
		t.Error("Registry eventConstructors map is nil")
	}
}

func TestRegisterEventType(t *testing.T) {
	reg := GetEventRegistry()

	constructor := func() Evaluable {
		return &MockEvaluable{}
	}

	reg.RegisterEventType("test_event", constructor)

	if _, exists := reg.eventConstructors["test_event"]; !exists {
		t.Error("Event type was not registered properly")
	}
}

func TestGetRegistry(t *testing.T) {
	reg := GetEventRegistry()
	constructor := func() Evaluable {
		return &MockEvaluable{}
	}

	reg.RegisterEventType("test_event", constructor)

	registry := reg.GetRegistry()
	if len(registry) == 0 {
		t.Error("Registry should not be empty")
	}
	if _, exists := registry["test_event"]; !exists {
		t.Error("Expected registered event type not found in registry")
	}
}

func TestProcessEvent(t *testing.T) {
	reg := GetEventRegistry()
	processor := &MockRuleProcessor{}

	tests := []struct {
		name        string
		eventJSON   string
		setup       func()
		wantSuccess bool
		wantErr     bool
	}{
		{
			name:        "Invalid JSON",
			eventJSON:   `{invalid json}`,
			setup:       func() {},
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:        "Missing event type",
			eventJSON:   `{"data": "test"}`,
			setup:       func() {},
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:        "Unregistered event type",
			eventJSON:   `{"type": "unknown_event"}`,
			setup:       func() {},
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:      "Valid event processing",
			eventJSON: `{"type": "test_event", "data": "test"}`,
			setup: func() {
				reg.RegisterEventType("test_event", func() Evaluable {
					return &MockEvaluable{
						EvaluateFunc: func(ctx context.Context, processor RuleProcessor) (bool, error) {
							return true, nil
						},
					}
				})
			},
			wantSuccess: true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			success, err := reg.ProcessEvent(context.Background(), processor, []byte(tt.eventJSON))

			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if success != tt.wantSuccess {
				t.Errorf("ProcessEvent() success = %v, want %v", success, tt.wantSuccess)
			}
		})
	}
}
