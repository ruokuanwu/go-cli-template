package db

import (
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
)

func TestOpen(t *testing.T) {
	if _, _, err := Open(sqlite.Open(filepath.Join(t.TempDir(), "missing", "test.db"))); err == nil {
		t.Fatal("expected invalid database path error")
	}
	database, cleanup, err := Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if err := database.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY)").Error; err != nil {
		t.Fatal(err)
	}
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Ping(); err == nil {
		t.Fatal("cleanup did not close database")
	}
}
