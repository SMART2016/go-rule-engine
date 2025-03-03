package models

import "time"

type Rule struct {
	RuleId        string        `json:"rule_id"`
	EventType     string        `json:"event_type"`
	Condition     string        `json:"condition"`     // Expression evaluated by grule
	Action        string        `json:"action"`        // Defines what to do when condition is met
	SendEmail     bool          `json:"send_email"`    // Whether to send an email
	Deduplication bool          `json:"deduplication"` // Whether to deduplicate events
	DedupWindow   time.Duration `json:"dedup_window"`  // Time window for deduplication (x hours)
	PayloadFields []string      `json:"payload_fields"`
}

// RuleSet represents a set of rules for a tenant.
type RuleSet map[string]Rule

// RuleConfig represents all rules for all tenants.
type RuleConfig map[string]RuleSet
