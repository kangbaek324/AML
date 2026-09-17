package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
	"github.com/kangbaek324/AML/internal/handler"
)

func New(queries *db.Queries, sourceQueries *sourcedb.Queries) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	h := handler.New(queries, sourceQueries)

	r.GET("/health", h.Health)

	api := r.Group("/api/v1")
	{
		api.GET("/users", h.ListUsers)
		api.GET("/users/:id", h.GetUser)
		api.GET("/accounts/:id", h.GetAccount)
	}

	return r
}
