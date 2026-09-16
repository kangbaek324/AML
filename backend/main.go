package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kangbaek324/AML/internal/config"
	"github.com/kangbaek324/AML/internal/router"
)

func main() {
	cfg := config.Load()

	conn, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	r := router.New()

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
