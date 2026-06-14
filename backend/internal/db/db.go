package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/internal/pkg/logger"

	_ "modernc.org/sqlite"
)

// Client is the Ent ORM client, exported for use by repositories.
var Client *ent.Client

// Init opens the database and initializes the Ent client.
func Init(driver, source string) error {
	dir := filepath.Dir(source)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	db, err := sql.Open(driver, source)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent performance
	if _, err = db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		logger.S().Warnw("Failed to enable WAL mode", "error", err)
	}

	// Enable foreign keys
	if _, err = db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		logger.S().Warnw("Failed to enable foreign keys", "error", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Use OpenDB to tell Ent this is a SQLite database
	drv := entsql.OpenDB(dialect.SQLite, db)
	Client = ent.NewClient(ent.Driver(drv))
	return nil
}

// Close closes the Ent client and underlying database connection.
func Close() {
	if Client != nil {
		_ = Client.Close()
	}
}

// Migrate runs Ent auto-migration to create/update schema.
func Migrate() error {
	ctx := context.Background()
	if err := Client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}
	logger.S().Infow("Database migration completed")
	return nil
}
