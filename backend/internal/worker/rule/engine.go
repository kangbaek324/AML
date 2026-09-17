package rule

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	db "github.com/kangbaek324/AML/db/sqlc"
)

// Engine은 등록된 Rule들을 각자의 Interval에 맞춰 독립된 goroutine으로 실행하고,
// 각 Rule의 워터마크(since)를 cursors 테이블에 영속화한다.
type Engine struct {
	queries *db.Queries
	rules   []Rule
}

func NewEngine(queries *db.Queries, rules ...Rule) *Engine {
	return &Engine{queries: queries, rules: rules}
}

func (e *Engine) Start(ctx context.Context) {
	for _, r := range e.rules {
		go e.runLoop(ctx, r)
	}
}

func (e *Engine) runLoop(ctx context.Context, r Rule) {
	since := e.loadCursor(ctx, r.Name())
	since = e.tick(ctx, r, since)

	ticker := time.NewTicker(r.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			since = e.tick(ctx, r, since)
		}
	}
}

func (e *Engine) tick(ctx context.Context, r Rule, since time.Time) time.Time {
	next, err := r.Run(ctx, since)
	if err != nil {
		log.Printf("rule %s: run failed: %v", r.Name(), err)
		return since
	}
	if next.Equal(since) {
		return since
	}

	if err := e.queries.UpsertCursor(ctx, db.UpsertCursorParams{
		Type:      r.Name(),
		Timestamp: next,
	}); err != nil {
		log.Printf("rule %s: save cursor failed: %v", r.Name(), err)
	}

	return next
}

func (e *Engine) loadCursor(ctx context.Context, ruleName string) time.Time {
	ts, err := e.queries.GetCursor(ctx, ruleName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("rule %s: load cursor failed: %v", ruleName, err)
		}
		return time.Now()
	}
	return ts
}
