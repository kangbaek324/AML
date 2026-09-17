package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
	"github.com/kangbaek324/AML/internal/handler"
	"github.com/kangbaek324/AML/internal/worker"
)

func New(queries *db.Queries, sourceQueries *sourcedb.Queries, w *worker.Worker) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	h := handler.New(queries, sourceQueries, w)

	r.GET("/health", h.Health)

	api := r.Group("/api/v1")
	{
		api.GET("/users", h.ListUsers)
		api.GET("/users/:id", h.GetUser)
		api.POST("/users/refresh", h.RefreshUsers)
		api.GET("/accounts/:id", h.GetAccount)
		api.GET("/alerts", h.ListAlerts)
		api.GET("/alerts/:id", h.GetAlert)
	}

	return r
}
