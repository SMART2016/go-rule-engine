package rule_processor

import (
	"context"
	"errors"
	"github.com/SMART2016/go-rule-engine/models"
	"github.com/SMART2016/go-rule-engine/store"
	"time"
)

type Config interface {
	GetDBConfigPath() string
	GetRuleRepoPath() string
	GetCleanupInterval() time.Duration
	DbConfig() *EventStateStoreConfig
}

type RuleProcessor[T any] interface {
	Evaluate(ctx context.Context, event models.Event[T]) (bool, error)
}

type EventProcessor[T any] struct {
	ruleProc RuleProcessor[T]
}

/*
*
This is the main API for the consumers of the framework to initialize the event processor.
This is a Go function named NewProcessor that creates a new instance of an event processor. It takes a configuration
object cfg as input and returns an EventProcessor instance and an error.
Here's a step-by-step breakdown:

It initializes a rule repository using the provided configuration.
It generates a database connection string (DSN) using the configuration's database settings.
It initializes an event store using the generated DSN.
It creates a new EventProcessor instance with a GRuleProcessor instance, which is configured with the initialized rule
repository and event store.
The function returns the EventProcessor instance and an error if any of the initialization steps fail.
*/
func NewProcessor[T any](cfg Config) (*EventProcessor[T], error) {
	//Initialize Rule Repository
	ruleRepo, err := initializeSingleRuleRepoInstance(cfg)
	if err != nil {
		return nil, errors.New("Failed to Initialize Rule Repository: " + err.Error())
	}

	//Generate DSN dsn := "host=localhost port=5432 user=username password=password dbname=mydb sslmode=disable"
	dsn := cfg.DbConfig().GenerateDSN()

	//Initialize Event Store
	eventStore, err := store.InitializeEventStateStore(dsn)
	if err != nil {
		return nil, errors.New("Failed to Initialize Event State Store: " + err.Error())
	}

	//TODO: Pass the Type as T
	return &EventProcessor[T]{
		ruleProc: &GRuleProcessor[T]{
			ruleRepo:   ruleRepo,
			eventStore: eventStore,
		},
	}, nil
}

/*
ProcessEvent evaluates the event against the configured rules.

It takes a context.Context and a models.Event[T] as an argument.
The context.Context is used to pass a context to the rule engine.
The models.Event[T] is the event to be evaluated.

It returns a boolean indicating whether the event should be handled and an error.
If the event should be handled, the boolean is true and the error is nil.
If the event should not be handled, the boolean is false and the error is nil.
If an error occurs, the boolean is false and the error is not nil.
*/
func (ep *EventProcessor[T]) ProcessEvent(ctx context.Context, event models.Event[T]) (bool, error) {
	// Evaluate the event using the rule engine
	return ep.ruleProc.Evaluate(ctx, event)
}
