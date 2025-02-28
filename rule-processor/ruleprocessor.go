package rule_processor

import (
	"context"
	"github.com/SMART2016/go-rule-engine/models"
	"github.com/SMART2016/go-rule-engine/store"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type RuleRepository interface {
	GetRules(tenantID, eventType string) ([]Rule, error)
}

type RuleProcessor[T any] struct {
	ruleRepo   RuleRepository
	eventStore store.Querier
}

func (re *RuleProcessor[T]) Evaluate(ctx context.Context, event models.Event[T]) (bool, error) {
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
		grl := `
rule ` + rule.Name + ` {
    when
        ` + rule.Condition + `
    then
        ` + rule.Action + `
}
`
		resource := pkg.NewBytesResource([]byte(grl))
		err := ruleBuilder.BuildRuleFromResource("EventRules", "0.0.1", resource)
		if err != nil {
			return false, err
		}

		// Retrieve the KnowledgeBase instance
		knowledgeBase, err := knowledgeLibrary.NewKnowledgeBaseInstance("EventRules", "0.0.1")

		// Create a new DataContext and add the event
		dataContext := ast.NewDataContext()
		err = dataContext.Add("Event", &event)
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

func NewRuleProcessor[T any](ruleRepo RuleRepository, eventStore store.Querier) *RuleProcessor[T] {
	return &RuleProcessor[T]{
		ruleRepo:   ruleRepo,
		eventStore: eventStore,
	}
}
