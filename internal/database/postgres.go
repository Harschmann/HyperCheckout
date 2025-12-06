package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	// CRITICAL: Connection pooling settings for high load

	// MaxOpenConns: Limits the number of open connections to the DB.
	// If 25 queries are running, the 26th will wait.
	// This prevents "too many connections" errors from postgres.
	db.SetMaxOpenConns(25)

	// MaxIdleConns: How many connections to keep open when no one is using the app.
	db.SetMaxIdleConns(25)

	// ConnMaxLifetime: Recycle connections every 5 minutes to prevent stale issues.
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("✅ Connected to Postgres successfully")
	return db, nil
}