package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	sourcedb "github.com/kangbaek324/AML/db/source/sqlc"
	db "github.com/kangbaek324/AML/db/sqlc"
	"github.com/kangbaek324/AML/internal/config"
	"github.com/kangbaek324/AML/internal/router"
	"github.com/kangbaek324/AML/internal/worker"
	"github.com/kangbaek324/AML/internal/worker/rule"
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

	sourceConn, err := sql.Open("mysql", cfg.SourceDSN())
	if err != nil {
		log.Fatalf("failed to open source db: %v", err)
	}
	defer sourceConn.Close()

	if err := sourceConn.Ping(); err != nil {
		log.Fatalf("failed to connect to source db: %v", err)
	}

	queries := db.New(conn)
	sourceQueries := sourcedb.New(sourceConn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userInfoWorker := worker.New(
		worker.NewAssetTierWorker(queries, sourceQueries),
		worker.NewUserRiskWorker(queries),
	)
	userInfoWorker.Start(ctx)

	ruleEngine := rule.NewEngine(
		queries,
		rule.NewLargeTransactionRule(queries, sourceQueries),
		// 새 룰은 여기에 추가하면 자동으로 각자의 주기로 실행된다.
	)
	ruleEngine.Start(ctx)

	r := router.New(queries, sourceQueries, userInfoWorker)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
