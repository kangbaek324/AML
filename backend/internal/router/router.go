package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kangbaek324/AML/internal/handler"
)

func New() *gin.Engine {
	r := gin.Default()

	h := handler.New()

	r.GET("/health", h.Health)

	return r
}
