package catcare

import (
	"context"
	"fmt"

	core "github.com/wastingnotime/zeroapps/core/catcare"
	"github.com/wastingnotime/zeroapps/store"
)

type CommandEnvelope struct {
	AggregateID     string
	Command         core.Command
	ExpectedVersion *int
}

type Result struct {
	Ok         bool
	NewVersion int
	Events     []core.Event
	Rejection  *core.Rejection
}

type Service struct {
	store      store.EventStore
	maxRetries int
}

func NewService(store store.EventStore) *Service {
	return &Service{store: store, maxRetries: 1}
}

func (s *Service) HandleCommand(ctx context.Context, env CommandEnvelope) (Result, error) {
	if env.AggregateID == "" {
		return Result{}, fmt.Errorf("aggregate id is required")
	}

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		rawEvents, version, err := s.store.Load(ctx, env.AggregateID)
		if err != nil {
			return Result{}, err
		}

		events, err := toCoreEvents(rawEvents)
		if err != nil {
			return Result{}, err
		}

		aggregate, err := core.LoadFrom(events)
		if err != nil {
			return Result{}, err
		}

		decided, err := aggregate.Decide(env.Command)
		if err != nil {
			if rejection, ok := err.(core.Rejection); ok {
				return Result{Ok: false, NewVersion: version, Rejection: &rejection}, nil
			}
			return Result{}, err
		}

		expected := version
		if env.ExpectedVersion != nil {
			expected = *env.ExpectedVersion
		}

		newVersion, err := s.store.Append(ctx, env.AggregateID, expected, toAnySlice(decided))
		if err == store.ErrConcurrencyConflict {
			if env.ExpectedVersion != nil {
				return Result{}, store.ErrConcurrencyConflict
			}
			if attempt < s.maxRetries {
				continue
			}
		}
		if err != nil {
			return Result{}, err
		}

		return Result{Ok: true, NewVersion: newVersion, Events: decided}, nil
	}

	return Result{}, store.ErrConcurrencyConflict
}

func toCoreEvents(events []any) ([]core.Event, error) {
	if len(events) == 0 {
		return nil, nil
	}

	typed := make([]core.Event, 0, len(events))
	for _, event := range events {
		typedEvent, ok := event.(core.Event)
		if !ok {
			return nil, fmt.Errorf("unexpected event type %T", event)
		}
		typed = append(typed, typedEvent)
	}
	return typed, nil
}

func toAnySlice(events []core.Event) []any {
	if len(events) == 0 {
		return nil
	}
	raw := make([]any, 0, len(events))
	for _, event := range events {
		raw = append(raw, event)
	}
	return raw
}
