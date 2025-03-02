package rule_processor

import (
	"context"
	"errors"
	"fmt"
	"github.com/SMART2016/go-rule-engine/models"
	"github.com/SMART2016/go-rule-engine/store"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type GRuleProcessor struct {
	conf       *Config
	ruleRepo   RuleRepository
	eventStore store.Querier
}

/*
NewGRuleProcessor initializes a new instance of GRuleProcessor.

Parameters:
  - cfg: *Config - The configuration settings for the rule processor.
  - ruleRepo: RuleRepository - Interface to access the rule repository.
  - eventStore: store.Querier - Interface to interact with the event store.

Returns:
  - *GRuleProcessor: A pointer to the initialized GRuleProcessor instance.
*/
func NewGRuleProcessor(cfg Config) (*GRuleProcessor, error) {
	// Initialize Rule Repository
	ruleRepo, err := initializeSingleRuleRepoInstance(cfg)
	if err != nil {
		return nil, errors.New("Failed to Initialize Rule Repository: " + err.Error())
	}

	// Generate DSN
	dsn := cfg.DbConfig().GenerateDSN()

	// Initialize BaseEvent Store
	eventStore, err := store.InitializeEventStateStore(dsn)
	if err != nil {
		return nil, errors.New("Failed to Initialize BaseEvent State Store: " + err.Error())
	}
	return &GRuleProcessor{
		conf:       &cfg,
		ruleRepo:   ruleRepo,
		eventStore: eventStore,
	}, nil
}

func (re *GRuleProcessor) Evaluate(ctx context.Context, event models.BaseEvent[any]) (bool, error) {
	err := event.Validate()
	if err != nil {
		return false, err
	}
	rules, err := re.ruleRepo.GetRules(event.TenantID, event.Type)
	if err != nil {
		return false, nil // No rules found for this tenant and event type
	}

	for _, rule := range rules {
		if rule.Deduplication {
			//TODO: populate store.IsDuplicateParams object from event and pass it to isDuplicate
			//generate SHA for the event based on the configured duplication key or atleast for now tenant_is,event_type and payload
			//store.IsDuplicateParams{event}
			isDuplicate, err := re.eventStore.IsDuplicate(ctx, store.IsDuplicateParams{})
			if err != nil {
				return false, err
			}
			if isDuplicate {
				continue // Skip processing for duplicate events
			}
		}

		// Initialize the KnowledgeLibrary and RuleBuilder
		knowledgeLibrary := ast.NewKnowledgeLibrary()
		ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

		// Build the rule into the KnowledgeBase
		grl := fmt.Sprintf(`
			rule %s {
				when
					%s
				then
					%s;
			}`, rule.Name, rule.Condition, rule.Action)
		resource := pkg.NewBytesResource([]byte(grl))
		err := ruleBuilder.BuildRuleFromResource("EventRules", "0.0.1", resource)
		if err != nil {
			return false, err
		}

		// Retrieve the KnowledgeBase instance
		knowledgeBase, err := knowledgeLibrary.NewKnowledgeBaseInstance("EventRules", "0.0.1")

		// Create a new DataContext and add the event
		dataContext := ast.NewDataContext()
		err = dataContext.Add("BaseEvent", &event)
		if err != nil {
			return false, err
		}

		// Initialize the Grule engine
		gruleEngine := engine.NewGruleEngine()

		// Execute the rules
		err = gruleEngine.Execute(dataContext, knowledgeBase)
		if err != nil {
			return false, err
		}

		// If the rule's action indicates an email should be sent
		if event.ShouldHandle {
			// Save the event to the store to track it for deduplication
			//TODO: populate store.SaveEventParams object from event result and pass it to SaveEvent
			err := re.eventStore.SaveEvent(ctx, store.SaveEventParams{event.TenantID, event.Type, event.EventSHA})
			if err != nil {
				return false, err
			}
			// Trigger email sending logic here
			return true, nil
		}
	}

	return false, nil
}
