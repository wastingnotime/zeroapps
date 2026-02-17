package store

import (
	"context"
	"sync"
)

type InMemoryStore struct {
	mu        sync.Mutex
	streams   map[string]*eventStream
	snapshots map[string]Snapshot
}

type eventStream struct {
	events  []any
	version int
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		streams:   map[string]*eventStream{},
		snapshots: map[string]Snapshot{},
	}
}

func (s *InMemoryStore) Load(ctx context.Context, streamID string) ([]any, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, exists := s.streams[streamID]
	if !exists {
		return nil, 0, nil
	}
	events := append([]any(nil), stream.events...)
	return events, stream.version, nil
}

func (s *InMemoryStore) Append(ctx context.Context, streamID string, expectedVersion int, events []any) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, exists := s.streams[streamID]
	if !exists {
		stream = &eventStream{}
		s.streams[streamID] = stream
	}

	if expectedVersion != stream.version {
		return stream.version, ErrConcurrencyConflict
	}

	if len(events) == 0 {
		return stream.version, nil
	}

	stream.events = append(stream.events, events...)
	stream.version += len(events)
	return stream.version, nil
}

func (s *InMemoryStore) LoadSnapshot(ctx context.Context, streamID string) (Snapshot, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot, ok := s.snapshots[streamID]
	return snapshot, ok, nil
}

func (s *InMemoryStore) SaveSnapshot(ctx context.Context, streamID string, snapshot Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.snapshots[streamID] = snapshot
	return nil
}

