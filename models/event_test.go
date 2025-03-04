package models

import (
	"context"
	"encoding/json"
	"testing"
)

type MockRuleProcessor2 struct {
	evaluateCalled bool
	returnValue    bool
	returnError    error
}

func (m *MockRuleProcessor2) Evaluate(ctx context.Context, event BaseEvent[any]) (bool, error) {
	m.evaluateCalled = true
	return m.returnValue, m.returnError
}

type TestPayload struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestBaseEvent_GenerateSHA256(t *testing.T) {
	event := BaseEvent[TestPayload]{
		TenantID: "tenant1",
		Type:     "test_event",
		Payload:  TestPayload{Field1: "test", Field2: 123},
	}

	tests := []struct {
		name      string
		dedupKeys string
		wantErr   bool
	}{
		{
			name:      "valid dedup keys",
			dedupKeys: "tenant1:test_event:123",
			wantErr:   false,
		},
		{
			name:      "empty dedup keys",
			dedupKeys: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sha, err := event.GenerateSHA256(tt.dedupKeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseEvent.GenerateSHA256() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(sha) == 0 {
				t.Error("BaseEvent.GenerateSHA256() returned empty SHA")
			}
		})
	}
}

func TestBaseEvent_ToJSON(t *testing.T) {
	tests := []struct {
		name    string
		event   BaseEvent[any]
		wantErr bool
	}{
		{
			name: "struct payload",
			event: BaseEvent[any]{
				Payload: TestPayload{Field1: "test", Field2: 123},
			},
			wantErr: false,
		},
		{
			name: "valid JSON string payload",
			event: BaseEvent[any]{
				Payload: `{"field1":"test","field2":123}`,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonStr, err := tt.event.ToJSON(tt.event.Payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseEvent.ToJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !json.Valid([]byte(jsonStr)) {
				t.Error("BaseEvent.ToJSON() returned invalid JSON")
			}
		})
	}
}

func TestBaseEvent_GetPayload(t *testing.T) {
	payload := TestPayload{Field1: "test", Field2: 123}
	event := BaseEvent[TestPayload]{
		Payload: payload,
	}

	result := event.GetPayload()
	if result != payload {
		t.Errorf("BaseEvent.GetPayload() = %v, want %v", result, payload)
	}
}
