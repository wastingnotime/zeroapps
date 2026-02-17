package store

import "context"

type EventStore interface {
	Load(ctx context.Context, streamID string) (events []any, version int, err error)
	Append(ctx context.Context, streamID string, expectedVersion int, events []any) (newVersion int, err error)
}

