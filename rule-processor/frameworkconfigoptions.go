package rule_processor

import (
	"time"
)

// Option defines a function signature for setting configurations.
type FrameworkConfigOption func(*FrameworkConfig)

// WithDBConfigPath sets the path to the database configuration file.
func WithDBConfigPath(path string) FrameworkConfigOption {
	return func(cfg *FrameworkConfig) {
		cfg.EventStoreConfigPath = path
	}
}

// WithRuleRepoPath sets the path to the rule repository JSON file.
func WithRuleRepoPath(path string) FrameworkConfigOption {
	return func(cfg *FrameworkConfig) {
		cfg.RuleRepoPath = path
	}
}

// WithCleanupInterval sets the event cleanup interval.
func WithCleanupInterval(interval time.Duration) FrameworkConfigOption {
	return func(cfg *FrameworkConfig) {
		cfg.CleanupInterval = interval
	}
}
