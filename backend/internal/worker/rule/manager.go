package rule

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	db "github.com/kangbaek324/AML/db/sqlc"
)

// Manager는 등록된 Rule들을 각자의 Interval에 맞춰 독립된 goroutine으로 실행하고,
// 각 Rule의 워터마크(since)를 cursors 테이블에 영속화한다.
type Manager struct {
	queries *db.Queries
	rules   []Rule
}

func NewManager(queries *db.Queries, rules ...Rule) *Manager {
	return &Manager{queries: queries, rules: rules}
}

func (m *Manager) Start(ctx context.Context) {
	for _, r := range m.rules {
		go m.runLoop(ctx, r)
	}
}

func (m *Manager) runLoop(ctx context.Context, r Rule) {
	since := m.loadCursor(ctx, r.Name())
	since = m.tick(ctx, r, since)

	ticker := time.NewTicker(r.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			since = m.tick(ctx, r, since)
		}
	}
}

func (m *Manager) tick(ctx context.Context, r Rule, since time.Time) time.Time {
	next, err := r.Run(ctx, since)
	if err != nil {
		log.Printf("rule %s: run failed: %v", r.Name(), err)
		return since
	}
	if next.Equal(since) {
		return since
	}

	if err := m.queries.UpsertCursor(ctx, db.UpsertCursorParams{
		Type:      r.Name(),
		Timestamp: next,
	}); err != nil {
		log.Printf("rule %s: save cursor failed: %v", r.Name(), err)
	}

	return next
}

func (m *Manager) loadCursor(ctx context.Context, ruleName string) time.Time {
	ts, err := m.queries.GetCursor(ctx, ruleName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("rule %s: load cursor failed: %v", ruleName, err)
		}
		return time.Now()
	}
	return ts
}
