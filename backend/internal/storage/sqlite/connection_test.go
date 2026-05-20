package sqlite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenAndRunMigrations(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	var tables []string
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		t.Fatalf("query tables: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table name: %v", err)
		}
		tables = append(tables, name)
	}

	expected := []string{"boards", "cards", "columns"}
	for _, et := range expected {
		found := false
		for _, tt := range tables {
			if tt == et {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected table %q not found in %v", et, tables)
		}
	}
}

func TestRunMigrationsIsIdempotent(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatalf("first migration run: %v", err)
	}

	if err := RunMigrations(db); err != nil {
		t.Fatalf("second migration run (idempotent): %v", err)
	}
}

func TestOpenSetsWALAndForeignKeys(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		t.Fatalf("query journal_mode: %v", err)
	}
	if journalMode != "wal" && journalMode != "WAL" {
		t.Fatalf("expected WAL journal mode, got %q", journalMode)
	}

	var foreignKeys int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		t.Fatalf("query foreign_keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("expected foreign_keys=1, got %d", foreignKeys)
	}
}

func TestOpenRejectsInvalidPath(t *testing.T) {
	t.Parallel()

	_, err := Open("/nonexistent/ directory/test.db")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
