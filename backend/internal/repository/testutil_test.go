package repository

import (
	_ "embed"
	"fmt"
	"testing"
)

// schemaSQL is embedded from testdata so tests can spin up an in-memory SQLite
// with the same schema as production. go:embed cannot traverse parent dirs,
// so a copy lives in ./testdata/schema.sql.
//
//go:embed testdata/schema.sql
var schemaSQL string

// newTestDB opens an in-memory SQLite database with the production schema applied.
// Each call returns a fresh database; modernc/sqlite honors ":memory:" per-connection
// when using a unique DSN, so we add a unique cache=shared name to isolate parallel tests.
func newTestDB(t *testing.T) *DB {
	t.Helper()
	dsn := fmt.Sprintf("file::memory:?cache=shared&_pragma=foreign_keys(1)")
	db, err := NewDB(dsn)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatalf("apply migration: %v", err)
	}
	// Reset tables to ensure isolation across parallel tests sharing the shared cache
	if _, err := db.Exec(`DELETE FROM weight_records; DELETE FROM health_records; DELETE FROM photos; DELETE FROM pets;`); err != nil {
		t.Fatalf("clean tables: %v", err)
	}
	return db
}
