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
	"log"
	"reflect"
)

type GRuleProcessor struct {
	conf       Config
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

	//TODO: Should instantiate Event Store ones and use it to get new connection and connect to it
	// Generate DSN
	//dsn := cfg.DbConfig().GenerateDSN()

	//// Initialize BaseEvent Store
	//eventStore, err := store.InitializeEventStateStore(dsn)
	//if err != nil {
	//	return nil, errors.New("Failed to Initialize BaseEvent State Store: " + err.Error())
	//}
	return &GRuleProcessor{
		conf:     cfg,
		ruleRepo: ruleRepo,
		//eventStore: eventStore,
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
		return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Event Validation failed %v", err)
	}

	rules, err := re.ruleRepo.GetRules(event.TenantID, event.Type)
	if err != nil {
		return false, nil // No rules found for this tenant and event type
	}
	// Generate DSN
	dsn := re.conf.DbConfig().GenerateDSN()

	// Initialize database
	conn, err := store.NewDatabase(dsn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
		return false, err
	}
	defer conn.Close() // Ensure closure of DB connection

	// Initialize Store
	eventStore := store.New(conn.DB)
	//b3f1b7856e68a4e004122b20c6c3bb43e312c50a3e25370362919f651b0c47a5
	for _, rule := range rules {
		if rule.Deduplication {
			//TODO: add rule_name or id too..
			isDuplicate, err := eventStore.IsDuplicate(ctx, store.IsDuplicateParams{
				TenantID:  event.TenantID,
				EventType: event.Type,
				EventSha:  event.EventSHA,
				Column5:   rule.DedupWindow,
				RuleID:    rule.RuleId,
			})
			if err != nil {
				return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Duplicate Event Check Failed: %v", err)
			}
			if isDuplicate {
				continue // Skip processing for duplicate events
			}
		}

		// Build the rule dynamically
		grl := fmt.Sprintf(`
			rule %s {
				when
					%s
				then
					%s;
			}
		`, rule.RuleId, rule.Condition, rule.Action)

		knowledgeLibrary := ast.NewKnowledgeLibrary()
		ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
		resource := pkg.NewBytesResource([]byte(grl))
		err = ruleBuilder.BuildRuleFromResource("EventRules", "0.0.1", resource)
		if err != nil {
			return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Rule Build Failed : %v", err)
		}

		// Retrieve the KnowledgeBase instance
		knowledgeBase, err := knowledgeLibrary.NewKnowledgeBaseInstance("EventRules", "0.0.1")
		if err != nil {
			return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Failed to get KnowledgeBase: %v", err)
		}

		// Create a new DataContext
		dataContext := ast.NewDataContext()

		err = dataContext.Add("Event", &event) // Add main event
		if err != nil {
			return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Failed to add Event to DataContext: %v", err)
		}

		// Extract and add Payload dynamically
		//payload := extractPayload(event.Payload)
		//if payload == nil {
		//	return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Failed to extract Payload")
		//}

		// Add Payload using interface
		payload := event.GetPayload()
		err = dataContext.Add("Payload", payload)
		if err != nil {
			return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Failed to add Payload to DataContext: %v", err)
		}

		err = dataContext.Add("Payload", payload)
		if err != nil {
			return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Failed to add Payload to DataContext: %v", err)
		}

		// Execute rules
		gruleEngine := engine.NewGruleEngine()
		err = gruleEngine.Execute(dataContext, knowledgeBase)
		if err != nil {
			return false, fmt.Errorf("[GRuleProcessor.Evaluate]: Rule Execution Failed : %v", err)
		}

		if event.ShouldHandle {
			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				return false, fmt.Errorf("Failed to convert payload to JSON: %v", err)
			}

			err = eventStore.SaveEvent(ctx, store.SaveEventParams{
				TenantID:  event.TenantID,
				EventType: event.Type,
				RuleID:    rule.RuleId,
				EventSha:  event.EventSHA,
				Column5:   json.RawMessage(jsonPayload),
			})
			if err != nil {
				return false, fmt.Errorf("Failed to save event to store: %v", err)
			}

			return true, nil
		}
	}

	return false, nil
}

/*
extractPayload extracts the payload from the event.

If the payload is a nil value, nil is returned.

If the payload is a pointer, the underlying value is returned.

If the payload is a struct, the struct is returned as an interface.

If the payload is a valid JSON string, it is unmarshaled into a map
and returned.

If the payload is not one of the above, nil is returned.
*/
func extractPayload(payload any) any {
	v := reflect.ValueOf(payload)

	// If payload is nil, return nil
	if !v.IsValid() || v.IsZero() {
		return nil
	}

	// If payload is a pointer, get the underlying value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// If payload is a struct, return it as an interface
	if v.Kind() == reflect.Struct {
		return payload
	}

	// If payload is a valid JSON string, unmarshal it into a map
	if str, ok := payload.(string); ok {
		var result map[string]interface{}
		err := json.Unmarshal([]byte(str), &result)
		if err == nil {
			return result
		}
	}

	return nil
}
