package store

import "context"

type Snapshot struct {
	Version int
	State   any
}

type SnapshotStore interface {
	LoadSnapshot(ctx context.Context, streamID string) (snapshot Snapshot, ok bool, err error)
	SaveSnapshot(ctx context.Context, streamID string, snapshot Snapshot) error
}

