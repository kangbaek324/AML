// Package migration은 db/migration 아래 SQL 파일들을 바이너리에 내장해,
// AML DB에 스키마 마이그레이션을 적용한다.
package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func Run(db *sql.DB) error {
	source, err := iofs.New(fs, ".")
	if err != nil {
		return fmt.Errorf("load migration source: %w", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("create mysql driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
