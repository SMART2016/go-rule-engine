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
	DbConfig() *DBConfig
}

type RuleProcessor[T any] interface {
	Evaluate(ctx context.Context, event models.Event[T]) (bool, error)
}

type EventProcessor[T any] struct {
	ruleProc RuleProcessor[T]
}

func NewEventProcessor[T any](cfg Config) (*EventProcessor[T], error) {
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

func (ep *EventProcessor[T]) ProcessEvent(ctx context.Context, event models.Event[T]) (bool, error) {
	// Evaluate the event using the rule engine
	return ep.ruleProc.Evaluate(ctx, event)
}
