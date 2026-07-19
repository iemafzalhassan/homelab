package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// DB wraps sql.DB with helper methods.
type DB struct {
	*sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection.
func NewPostgresDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.DB.Close()
}

// Migrate runs all pending migrations.
func (db *DB) Migrate() error {
	// Use goose to run migrations.
	// The migration files are located in the kyros/database/migrations directory.
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db.DB, "kyros/database/migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}