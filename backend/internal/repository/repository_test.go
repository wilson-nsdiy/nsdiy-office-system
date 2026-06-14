package repository

import (
	"context"
	"database/sql"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"oa-nsdiy/backend/ent"

	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite database with Ent for testing.
func setupTestDB(t *testing.T) *ent.Client {
	t.Helper()

	db, err := sql.Open("sqlite", "file:test.db?mode=memory&cache=shared&_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))

	// Run auto-migration
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		t.Fatalf("failed to migrate test db: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	return client
}
