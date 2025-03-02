package rule_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// RuleRepository provides access to rules.
type singletonJsonRuleRepository struct {
	cfg   Config
	rules map[string]map[string][]Rule
	mu    sync.RWMutex
}

var (
	instance *singletonJsonRuleRepository
	once     sync.Once
)

// GetInstance returns the singleton instance of Singleton
func initializeSingleRuleRepoInstance(frameWrkCfg Config) (*singletonJsonRuleRepository, error) {
	var instantiationErr error = nil
	once.Do(func() {
		file, err := os.ReadFile(frameWrkCfg.GetRuleRepoPath())
		if err != nil {
			instantiationErr = err
			return
		}
		if err = validateRules("rules.json"); err != nil {
			instantiationErr = errors.New(fmt.Sprintf("Rules Validation failed: %+v", err))
		}
		var r map[string]map[string][]Rule
		err = json.Unmarshal(file, &r)
		if err != nil {
			instantiationErr = err
			return
		}

		instance = &singletonJsonRuleRepository{
			rules: r,
			cfg:   frameWrkCfg,
		}
	})
	return instance, instantiationErr
}

// GetRules retrieves rules for a specific tenant and event type.
func (r *singletonJsonRuleRepository) GetRules(tenantID, eventType string) ([]Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if tenantRules, ok := r.rules[tenantID]; ok {
		if eventRules, ok := tenantRules[eventType]; ok {
			return eventRules, nil
		}
	}
	return nil, nil
}
