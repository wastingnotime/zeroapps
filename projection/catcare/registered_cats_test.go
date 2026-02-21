package catcare

import (
	"context"
	"testing"

	core "github.com/wastingnotime/zeroapps/core/catcare"
)

func TestListRegisteredCatsGivenEmptyProjectionWhenListThenReturnsEmpty(t *testing.T) {
	projection := NewRegisteredCats()

	cats := projection.ListRegisteredCats()
	if len(cats) != 0 {
		t.Fatalf("expected empty list, got %d cats", len(cats))
	}
}

func TestListRegisteredCatsGivenSingleRegisteredCatWhenListThenReturnsCat(t *testing.T) {
	projection := NewRegisteredCats()

	err := projection.Apply(context.Background(), "cat-1", 1, core.CatRegistered{
		CommandID: "cmd-1",
		CatID:     "cat-1",
		Name:      "Miso",
		BirthDate: "2023-01-01",
	})
	if err != nil {
		t.Fatalf("apply event: %v", err)
	}

	cats := projection.ListRegisteredCats()
	if len(cats) != 1 {
		t.Fatalf("expected 1 cat, got %d", len(cats))
	}
	if cats[0].CatID != "cat-1" || cats[0].Name != "Miso" || cats[0].BirthDate != "2023-01-01" {
		t.Fatalf("unexpected cat %+v", cats[0])
	}
}

func TestListRegisteredCatsGivenMultipleRegisteredCatsWhenListThenReturnsAllSortedByCatID(t *testing.T) {
	projection := NewRegisteredCats()

	events := []struct {
		streamID string
		version  int
		event    core.Event
	}{
		{
			streamID: "cat-2",
			version:  1,
			event: core.CatRegistered{
				CommandID: "cmd-2",
				CatID:     "cat-2",
				Name:      "Taro",
				BirthDate: "2023-05-02",
			},
		},
		{
			streamID: "cat-1",
			version:  1,
			event: core.CatRegistered{
				CommandID: "cmd-1",
				CatID:     "cat-1",
				Name:      "Miso",
				BirthDate: "2023-01-01",
			},
		},
	}

	for _, item := range events {
		if err := projection.Apply(context.Background(), item.streamID, item.version, item.event); err != nil {
			t.Fatalf("apply event: %v", err)
		}
	}

	cats := projection.ListRegisteredCats()
	if len(cats) != 2 {
		t.Fatalf("expected 2 cats, got %d", len(cats))
	}
	if cats[0].CatID != "cat-1" || cats[1].CatID != "cat-2" {
		t.Fatalf("expected sorted cats by id, got %+v", cats)
	}
}

func TestListRegisteredCatsGivenReappliedVersionWhenApplyThenDoesNotDuplicate(t *testing.T) {
	projection := NewRegisteredCats()
	registered := core.CatRegistered{
		CommandID: "cmd-1",
		CatID:     "cat-1",
		Name:      "Miso",
		BirthDate: "2023-01-01",
	}

	if err := projection.Apply(context.Background(), "cat-1", 1, registered); err != nil {
		t.Fatalf("apply event: %v", err)
	}
	if err := projection.Apply(context.Background(), "cat-1", 1, registered); err != nil {
		t.Fatalf("apply event again: %v", err)
	}

	cats := projection.ListRegisteredCats()
	if len(cats) != 1 {
		t.Fatalf("expected exactly 1 cat, got %d", len(cats))
	}
}
