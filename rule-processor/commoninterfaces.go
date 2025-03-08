package rule_processor

import (
	"context"
	"github.com/SMART2016/go-rule-engine/models"
	"github.com/SMART2016/go-rule-engine/store"
	"time"
)

/*
RuleProcessor is the interface that must be implemented by a rule processor.

It provides a single function, Evaluate, which takes a context and a
BaseEvent and returns a boolean indicating whether the event was handled
and an error if there was a problem.
*/
type RuleProcessor interface {
	// Evaluate takes a context and a BaseEvent and returns a boolean indicating
	// whether the event was handled and an error if there was a problem.
	Evaluate(ctx context.Context, event models.BaseEvent[any]) (bool, error)
}

/*
RuleRepository is an interface that provides access to rules.

The GetRules function retrieves all rules for a given tenant and event
type.
*/
type RuleRepository interface {
	// GetRules retrieves all rules for a given tenant and event type.
	//
	// It takes a tenant ID and an event type as input and returns a slice of
	// Rule objects and an error if there was a problem.
	GetRules(tenantID, eventType string) ([]models.Rule, error)
}

/*
Config is the interface that must be implemented by a configuration
provider.

A configuration provider is responsible for providing the information
necessary to configure the rule processor.

The information provided by the configuration provider includes the path
to the database configuration file, the path to the rule repository JSON
file, the cleanup interval, and the database configuration.
*/
type Config interface {
	// GetDBConfigPath returns the path to the database configuration file.
	GetDBConfigPath() string

	// GetRuleRepoPath returns the path to the rule repository JSON file.
	GetRuleRepoPath() string

	// GetCleanupInterval returns the cleanup interval.
	GetCleanupInterval() time.Duration

	// DbConfig returns the database configuration.
	DbConfig() *EventStateStoreConfig
}

/*
EventStore is an interface for managing events.

It provides methods to check for duplicate events, save new events, and
cleanup old events from the database.
*/
type EventStore interface {
	// IsDuplicate checks if an event is a duplicate.
	//
	// It takes a context, a database transaction, and parameters for checking
	// duplication, returning a boolean indicating duplication status and an error if
	// there was a problem.
	IsDuplicate(ctx context.Context, db store.DBTX, arg store.IsDuplicateParams) (interface{}, error)
	SaveEvent(ctx context.Context, db store.DBTX, arg store.SaveEventParams) error
	CleanupOldEvents(ctx context.Context, db store.DBTX, dollar_1 interface{}) error
}
