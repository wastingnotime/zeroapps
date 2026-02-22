package store

import (
	"context"
	"path/filepath"
	"testing"

	core "github.com/wastingnotime/zeroapps/core/catcare"
)

func TestSQLiteStoreGivenAppendedEventsWhenLoadThenReturnsEventsAndVersion(t *testing.T) {
	ctx := context.Background()
	store := newSQLiteStoreForTest(t)
	t.Cleanup(func() {
		_ = store.Close()
	})

	newVersion, err := store.Append(ctx, "cat-cmd-1", 0, []any{core.CatRegistered{
		CommandID: "cmd-1",
		CatID:     "cat-cmd-1",
		Name:      "Miso",
		BirthDate: "2023-01-01",
	}})
	if err != nil {
		t.Fatalf("append: %v", err)
	}
	if newVersion != 1 {
		t.Fatalf("newVersion = %d, want 1", newVersion)
	}

	events, version, err := store.Load(ctx, "cat-cmd-1")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if version != 1 {
		t.Fatalf("version = %d, want 1", version)
	}
	if len(events) != 1 {
		t.Fatalf("len(events) = %d, want 1", len(events))
	}
	registered, ok := events[0].(core.CatRegistered)
	if !ok {
		t.Fatalf("event type = %T, want core.CatRegistered", events[0])
	}
	if registered.CatID != "cat-cmd-1" {
		t.Fatalf("cat_id = %q, want cat-cmd-1", registered.CatID)
	}
}

func TestSQLiteStoreGivenWrongExpectedVersionWhenAppendThenReturnsConcurrencyConflict(t *testing.T) {
	ctx := context.Background()
	store := newSQLiteStoreForTest(t)
	t.Cleanup(func() {
		_ = store.Close()
	})

	_, err := store.Append(ctx, "cat-cmd-1", 0, []any{core.CatRegistered{
		CommandID: "cmd-1",
		CatID:     "cat-cmd-1",
		Name:      "Miso",
		BirthDate: "2023-01-01",
	}})
	if err != nil {
		t.Fatalf("initial append: %v", err)
	}

	_, err = store.Append(ctx, "cat-cmd-1", 0, []any{core.WeightLogged{
		CommandID: "cmd-2",
		EntryID:   "weight-cmd-2",
		At:        "2026-02-14T10:00:00Z",
		Grams:     4200,
	}})
	if err != ErrConcurrencyConflict {
		t.Fatalf("err = %v, want %v", err, ErrConcurrencyConflict)
	}
}

func newSQLiteStoreForTest(t *testing.T) *SQLiteStore {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "catcare.db")
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	return store
}
