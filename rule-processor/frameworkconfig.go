package rule_processor

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// DBConfig holds database connection details for the State Store.
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

// FrameworkConfig holds all configuration for the rule engine.
type FrameworkConfig struct {
	DBConfigPath    string
	RuleRepoPath    string
	CleanupInterval time.Duration
	dbConfig        *DBConfig
	rules           map[string]map[string][]Rule
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

// NewFrameworkConfig initializes a new configuration with functional options.
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
	file, err := os.ReadFile(cfg.DBConfigPath)
	if err != nil {
		return err
	}
	var dbConfig DBConfig
	err = json.Unmarshal(file, &dbConfig)
	if err != nil {
		return err
	}
	cfg.dbConfig = &dbConfig
	return nil
}
