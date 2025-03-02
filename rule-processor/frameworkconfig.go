/*
*
Users will create the configuration and pass it to the event processor for evaluating the event against the configured rules.
This module provides configuration management for a rule processing framework. It includes event state store configurations details,
rule repository management,and configurable cleanup intervals. The configuration is loaded dynamically from files
provided by the user.
*/
package rule_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// EventStateStoreConfig holds database connection details for the State Store.
type EventStateStoreConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

// GenerateDSN constructs a PostgreSQL DSN from EventStateStoreConfig
func (cfg *EventStateStoreConfig) GenerateDSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)
}

// FrameworkConfig holds all configuration for the rule engine.
type FrameworkConfig struct {
	EventStoreConfigPath string
	RuleRepoPath         string
	CleanupInterval      time.Duration
	eventStoreConfig     *EventStateStoreConfig
	rules                map[string]map[string][]Rule
}

/*
NewFrameworkConfig initializes a new configuration with functional options.

The following options are available:

- WithDBConfigPath(string): sets the path to the database configuration file.

- WithRuleRepoPath(string): sets the path to the rule repository JSON file.

- WithCleanupInterval(time.Duration): sets the event cleanup interval.

The provided options are applied to the configuration in order. If an option
is not provided, the default value is used.

The configurations from the provided paths are loaded after all options are
applied. If the loading process fails, an error is returned.
*/
func NewFrameworkConfig(opts ...FrameworkConfigOption) (*FrameworkConfig, error) {
	cfg := &FrameworkConfig{
		CleanupInterval: 24 * time.Hour, // Default cleanup interval
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// Load configurations from provided paths
	if err := cfg.Load(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *FrameworkConfig) GetDBConfigPath() string {
	return cfg.EventStoreConfigPath
}

func (cfg *FrameworkConfig) GetRuleRepoPath() string {
	return cfg.RuleRepoPath
}

func (cfg *FrameworkConfig) GetCleanupInterval() time.Duration {
	return cfg.CleanupInterval
}

func (cfg *FrameworkConfig) DbConfig() *EventStateStoreConfig {
	return cfg.eventStoreConfig
}

// Rule represents a single rule in the rule repository.
type Rule struct {
	Name          string        `json:"name"`
	EventType     string        `json:"event_type"`
	Condition     string        `json:"condition"`
	Action        string        `json:"action"`
	SendEmail     bool          `json:"send_email"`
	Deduplication bool          `json:"deduplication"`
	DedupWindow   time.Duration `json:"dedup_window"`
}

func (cfg *FrameworkConfig) Load() error {
	//Load DB config from the provided path by the consumer.
	if err := cfg.LoadDBConfig(); err != nil {
		return errors.New("load db config failed, Error : " + err.Error())
	}

	//Initialize a Rule repository instance and load the rules to the repository
	if _, err := initializeSingleRuleRepoInstance(cfg); err != nil {
		return errors.New("Initializing Rules Repository failed, Error : " + err.Error())
	}
	return nil
}

// LoadDBConfig loads the database configuration from a JSON file.
func (cfg *FrameworkConfig) LoadDBConfig() error {
	file, err := os.ReadFile(cfg.EventStoreConfigPath)
	if err != nil {
		return err
	}
	var dbConfig EventStateStoreConfig
	err = json.Unmarshal(file, &dbConfig)
	if err != nil {
		return err
	}
	cfg.eventStoreConfig = &dbConfig
	return nil
}
