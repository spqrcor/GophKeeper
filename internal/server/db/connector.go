// Package db работа с database
package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

// Connect соединение с db, DatabaseDSN - параметры подключения
func Connect(DatabaseDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
