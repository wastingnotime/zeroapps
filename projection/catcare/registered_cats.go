package catcare

import (
	"context"
	"sort"
	"sync"

	core "github.com/wastingnotime/zeroapps/core/catcare"
)

type RegisteredCat struct {
	CatID     string
	Name      string
	BirthDate string
}

type RegisteredCats struct {
	mu                sync.RWMutex
	catsByID          map[string]RegisteredCat
	lastStreamVersion map[string]int
}

func NewRegisteredCats() *RegisteredCats {
	return &RegisteredCats{
		catsByID:          map[string]RegisteredCat{},
		lastStreamVersion: map[string]int{},
	}
}

func (p *RegisteredCats) Apply(_ context.Context, streamID string, version int, event core.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	lastVersion := p.lastStreamVersion[streamID]
	if version <= lastVersion {
		return nil
	}

	switch ev := event.(type) {
	case core.CatRegistered:
		p.catsByID[ev.CatID] = RegisteredCat{
			CatID:     ev.CatID,
			Name:      ev.Name,
			BirthDate: ev.BirthDate,
		}
	}

	p.lastStreamVersion[streamID] = version
	return nil
}

func (p *RegisteredCats) ListRegisteredCats() []RegisteredCat {
	p.mu.RLock()
	defer p.mu.RUnlock()

	cats := make([]RegisteredCat, 0, len(p.catsByID))
	for _, cat := range p.catsByID {
		cats = append(cats, cat)
	}
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].CatID < cats[j].CatID
	})
	return cats
}
