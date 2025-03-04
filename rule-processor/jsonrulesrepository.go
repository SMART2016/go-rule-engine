package rule_processor

import (
	"encoding/json"
	"github.com/SMART2016/go-rule-engine/models"
	"os"
	"sync"
)

const (
	DEFAULT_TENANT_RULET_ID = "tenant_default"
)

// RuleRepository provides access to rules.
type singletonJsonRuleRepository struct {
	cfg   Config
	rules map[string][]models.Rule
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
		//TODO: Handle rules validation against the event schema seperately
		//if err = validateRules("rules.json"); err != nil {
		//	instantiationErr = errors.New(fmt.Sprintf("Rules Validation failed: %+v", err))
		//}
		var r map[string][]models.Rule
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

/*
GetRules retrieves rules for a specific tenant and event type.
GetRules retrieves all rules for a specified tenant and event type.

It acquires a read lock to ensure thread-safe access to the rules map.
If the tenant's rules exist, it filters rules matching the given event type.
It returns a slice of matching rules or an empty slice if none are found.
An error is returned if any issues occur during retrieval.
*/
func (r *singletonJsonRuleRepository) GetRules(tenantID, eventType string) ([]models.Rule, error) {
	r.mu.RLock()         // Acquire read lock for thread-safe access
	defer r.mu.RUnlock() // Ensure lock is released

	var rules []models.Rule
	if tenantRules, ok := r.rules[tenantID]; ok {
		for _, rule := range tenantRules {
			if eventType == rule.EventType {
				rules = append(rules, rule)
			}
		}
	} else {
		if tenantRules, ok = r.rules[DEFAULT_TENANT_RULET_ID]; ok {
			for _, rule := range tenantRules {
				if eventType == rule.EventType {
					rules = append(rules, rule)
				}
			}
		}
	}

	return rules, nil
}
