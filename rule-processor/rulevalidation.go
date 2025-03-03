package rule_processor

import (
	"encoding/json"
	"fmt"
	"github.com/SMART2016/go-rule-engine/examples/events"

	"github.com/SMART2016/go-rule-engine/models"
	"os"
	"reflect"
)

// Map event types to their corresponding payload struct.
var eventPayloadTypes map[string]reflect.Type = map[string]reflect.Type{
	"disk_space": reflect.TypeOf(events.DiskUsagePayload{}),
}

// ValidateRules loads rules from a JSON file and validates their `payload_fields`.
func validateRules(filePath string) error {
	// Step 1: Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read rules file: %w", err)
	}

	// Step 2: Unmarshal JSON into RuleConfig
	var ruleConfig models.RuleConfig
	if err := json.Unmarshal(data, &ruleConfig); err != nil {
		return fmt.Errorf("failed to parse rules JSON: %w", err)
	}

	// Step 3: Validate each rule
	for tenantID, ruleSet := range ruleConfig {
		for ruleID, rule := range ruleSet {
			if err := validateRulePayloadFields(rule); err != nil {
				return fmt.Errorf("tenant: %s, rule: %s, error: %w", tenantID, ruleID, err)
			}
		}
	}

	fmt.Println("All rules validated successfully.")
	return nil
}

// validateRulePayloadFields checks if all `payload_fields` exist in the event's payload struct.
func validateRulePayloadFields(rule models.Rule) error {
	// Step 1: Get the expected struct type for this event type
	payloadType, exists := eventPayloadTypes[rule.EventType]
	if !exists {
		return fmt.Errorf("unknown event type: %s", rule.EventType)
	}

	// Step 2: Validate each field in `payload_fields`
	for _, field := range rule.PayloadFields {
		if _, found := payloadType.FieldByName(field); !found {
			return fmt.Errorf("invalid field '%s' in payload_fields for event type '%s'", field, rule.EventType)
		}
	}

	return nil
}
