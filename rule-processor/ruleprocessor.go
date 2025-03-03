package rule_processor

import (
	"context"
	"encoding/json"
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

/*
Evaluate takes a context and a BaseEvent and returns a boolean indicating
whether the event was handled and an error if there was a problem.

The method first validates the event. If the event is invalid, an error is
returned. If the event is valid, it retrieves the rules associated with the
event's tenant and type from the rule repository. If no rules are found, the
method returns false, nil. If rules are found, the method iterates over the
rules and for each rule:

1. If the rule has deduplication enabled, the method checks if the event is a
duplicate by checking the event's SHA in the event store. If the event is a
duplicate, the method skips processing the rule.

2. The method builds the rule using the RuleBuilder and adds it to the
KnowledgeBase.

3. The method creates a new DataContext and adds the event to it.

4. The method initializes the Grule engine and executes the rules in the
KnowledgeBase.

5. If the rule's action indicates an email should be sent, the method saves the
event to the event store and triggers the email sending logic.

Parameters:
  - ctx: context.Context - A context to manage cancellation and deadlines.
  - event: models.BaseEvent[any] - The event to be evaluated.

Returns:
  - bool - Indicates whether the event was handled successfully.
  - error - Contains any error encountered during processing or evaluation of

the event.
*/
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
			//Generate SHA for the event based on the configured duplication key or atleast for now tenant_id,event_type and payload
			//store.IsDuplicateParams{event}
			isDuplicate, err := re.eventStore.IsDuplicate(ctx, store.IsDuplicateParams{
				TenantID:  event.TenantID,
				EventType: event.Type,
				EventSha:  event.EventSHA,
				Column4:   rule.DedupWindow,
			})
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
			}
`, rule.Name, rule.Condition, rule.Action)
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
			jsonPayload, err := event.ToJSON(event.Payload)
			if err != nil {
				return false, errors.New(fmt.Sprintf("Failed to convert payload to JSON: %v", err.Error()))
			}
			err = re.eventStore.SaveEvent(ctx, store.SaveEventParams{
				TenantID:  event.TenantID,
				EventType: event.Type,
				EventSha:  event.EventSHA,
				Column4:   json.RawMessage(jsonPayload),
			})
			if err != nil {
				return false, errors.New(fmt.Sprintf("Failed to save event to store: %v", err.Error()))
			}
			// Trigger email sending logic here
			return true, nil
		}
	}

	return false, nil
}
