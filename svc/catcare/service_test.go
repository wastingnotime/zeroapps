package catcare

import (
	"context"
	"testing"

	core "github.com/wastingnotime/zeroapps/core/catcare"
	"github.com/wastingnotime/zeroapps/store"
)

type projectionCall struct {
	streamID string
	version  int
	event    core.Event
}

type spyProjector struct {
	calls []projectionCall
}

func (p *spyProjector) Apply(_ context.Context, streamID string, version int, event core.Event) error {
	p.calls = append(p.calls, projectionCall{
		streamID: streamID,
		version:  version,
		event:    event,
	})
	return nil
}

func TestHandleCommandHappyPath(t *testing.T) {
	eventStore := store.NewInMemoryStore()
	projection := &spyProjector{}
	service := NewService(eventStore, projection)

	result, err := service.HandleCommand(context.Background(), CommandEnvelope{
		AggregateID: "cat-1",
		Command: core.RegisterCat{
			CommandID: "cmd-1",
			Name:      "Miso",
			BirthDate: "2023-01-01",
		},
	})
	if err != nil {
		t.Fatalf("handle command: %v", err)
	}
	if !result.Ok {
		t.Fatalf("expected ok result, got rejection %v", result.Rejection)
	}
	if result.NewVersion != 1 {
		t.Fatalf("expected version 1, got %d", result.NewVersion)
	}
	if len(result.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result.Events))
	}
	if len(projection.calls) != 1 {
		t.Fatalf("expected 1 projection call, got %d", len(projection.calls))
	}
	if projection.calls[0].streamID != "cat-1" {
		t.Fatalf("expected stream cat-1, got %q", projection.calls[0].streamID)
	}
	if projection.calls[0].version != 1 {
		t.Fatalf("expected version 1, got %d", projection.calls[0].version)
	}
}

func TestHandleCommandExpectedVersionConflict(t *testing.T) {
	eventStore := store.NewInMemoryStore()
	service := NewService(eventStore)

	_, err := service.HandleCommand(context.Background(), CommandEnvelope{
		AggregateID: "cat-1",
		Command: core.RegisterCat{
			CommandID: "cmd-1",
			Name:      "Miso",
		},
	})
	if err != nil {
		t.Fatalf("seed register: %v", err)
	}

	stale := 0
	_, err = service.HandleCommand(context.Background(), CommandEnvelope{
		AggregateID:     "cat-1",
		ExpectedVersion: &stale,
		Command: core.LogWeight{
			CommandID: "cmd-2",
			At:        "2026-02-14T10:00:00Z",
			Grams:     4200,
		},
	})
	if err != store.ErrConcurrencyConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
}
