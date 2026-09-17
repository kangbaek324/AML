package handler

import (
	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
	"github.com/kangbaek324/AML/internal/worker/userinfo"
)

type Handler struct {
	Queries       *db.Queries
	SourceQueries *sourcedb.Queries
	Worker        *userinfo.Manager
}

func New(queries *db.Queries, sourceQueries *sourcedb.Queries, w *userinfo.Manager) *Handler {
	return &Handler{Queries: queries, SourceQueries: sourceQueries, Worker: w}
}
