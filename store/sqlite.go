package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	core "github.com/wastingnotime/zeroapps/core/catcare"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

type sqliteEventRow struct {
	Version int
	Type    string
	Payload string
}

type sqliteProjectionApplier interface {
	Apply(ctx context.Context, streamID string, version int, event core.Event) error
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("db path is required")
	}

	if dbPath != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.initSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SQLiteStore) initSchema(ctx context.Context) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS streams (
	stream_id TEXT PRIMARY KEY,
	version INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
	stream_id TEXT NOT NULL,
	version INTEGER NOT NULL,
	event_type TEXT NOT NULL,
	payload TEXT NOT NULL,
	PRIMARY KEY (stream_id, version)
);

CREATE INDEX IF NOT EXISTS idx_events_stream_version
ON events(stream_id, version);
`
	_, err := s.db.ExecContext(ctx, ddl)
	return err
}

func (s *SQLiteStore) Load(ctx context.Context, streamID string) ([]any, int, error) {
	version, err := s.streamVersion(ctx, streamID)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT version, event_type, payload
FROM events
WHERE stream_id = ?
ORDER BY version ASC
`, streamID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	events := make([]any, 0)
	for rows.Next() {
		var row sqliteEventRow
		if err := rows.Scan(&row.Version, &row.Type, &row.Payload); err != nil {
			return nil, 0, err
		}
		event, err := decodeCatCareEvent(row.Type, row.Payload)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return events, version, nil
}

func (s *SQLiteStore) Append(ctx context.Context, streamID string, expectedVersion int, events []any) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	currentVersion, err := s.streamVersionTx(ctx, tx, streamID)
	if err != nil {
		return 0, err
	}
	if expectedVersion != currentVersion {
		return currentVersion, ErrConcurrencyConflict
	}
	if len(events) == 0 {
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return currentVersion, nil
	}

	for index, rawEvent := range events {
		event, ok := rawEvent.(core.Event)
		if !ok {
			return 0, fmt.Errorf("unexpected event type %T", rawEvent)
		}
		eventType, payload, err := encodeCatCareEvent(event)
		if err != nil {
			return 0, err
		}
		eventVersion := currentVersion + index + 1
		if _, err := tx.ExecContext(ctx, `
INSERT INTO events(stream_id, version, event_type, payload)
VALUES(?, ?, ?, ?)
`, streamID, eventVersion, eventType, payload); err != nil {
			return 0, err
		}
	}

	newVersion := currentVersion + len(events)
	if _, err := tx.ExecContext(ctx, `
INSERT INTO streams(stream_id, version) VALUES(?, ?)
ON CONFLICT(stream_id) DO UPDATE SET version = excluded.version
`, streamID, newVersion); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newVersion, nil
}

func (s *SQLiteStore) Replay(ctx context.Context, projector sqliteProjectionApplier) error {
	rows, err := s.db.QueryContext(ctx, `
SELECT stream_id, version, event_type, payload
FROM events
ORDER BY stream_id ASC, version ASC
`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var streamID string
		var version int
		var eventType string
		var payload string
		if err := rows.Scan(&streamID, &version, &eventType, &payload); err != nil {
			return err
		}
		event, err := decodeCatCareEvent(eventType, payload)
		if err != nil {
			return err
		}
		if err := projector.Apply(ctx, streamID, version, event); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (s *SQLiteStore) streamVersion(ctx context.Context, streamID string) (int, error) {
	var version int
	err := s.db.QueryRowContext(ctx, `
SELECT version FROM streams WHERE stream_id = ?
`, streamID).Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return version, err
}

func (s *SQLiteStore) streamVersionTx(ctx context.Context, tx *sql.Tx, streamID string) (int, error) {
	var version int
	err := tx.QueryRowContext(ctx, `
SELECT version FROM streams WHERE stream_id = ?
`, streamID).Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return version, err
}

func encodeCatCareEvent(event core.Event) (string, string, error) {
	switch ev := event.(type) {
	case core.CatRegistered:
		payload, err := json.Marshal(ev)
		return "CatRegistered", string(payload), err
	case core.WeightLogged:
		payload, err := json.Marshal(ev)
		return "WeightLogged", string(payload), err
	default:
		return "", "", fmt.Errorf("unsupported event type %T", event)
	}
}

func decodeCatCareEvent(eventType string, payload string) (core.Event, error) {
	switch eventType {
	case "CatRegistered":
		var event core.CatRegistered
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "WeightLogged":
		var event core.WeightLogged
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return nil, err
		}
		return event, nil
	default:
		return nil, fmt.Errorf("unsupported event type %q", eventType)
	}
}
