package db

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"os"
)

const path = "internal/server/db/migrations"

// Migrate запуск миграций, db - ресурс подключение к db
func Migrate(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	_, err := os.Stat(path)
	if err == nil {
		if err := goose.Up(db, path); err != nil {
			return err
		}
	}
	return nil
}
