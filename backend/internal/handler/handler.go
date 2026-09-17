package handler

import (
	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
)

type Handler struct {
	Queries       *db.Queries
	SourceQueries *sourcedb.Queries
}

func New(queries *db.Queries, sourceQueries *sourcedb.Queries) *Handler {
	return &Handler{Queries: queries, SourceQueries: sourceQueries}
}
