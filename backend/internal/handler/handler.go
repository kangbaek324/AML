package handler

import (
	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
	"github.com/kangbaek324/AML/internal/worker"
)

type Handler struct {
	Queries       *db.Queries
	SourceQueries *sourcedb.Queries
	Worker        *worker.Worker
}

func New(queries *db.Queries, sourceQueries *sourcedb.Queries, w *worker.Worker) *Handler {
	return &Handler{Queries: queries, SourceQueries: sourceQueries, Worker: w}
}
