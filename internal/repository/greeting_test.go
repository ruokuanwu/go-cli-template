package repository

import (
	"path/filepath"
	"testing"

	"example.com/go-cli-template/internal/db"
	"example.com/go-cli-template/internal/model"
	"github.com/glebarez/sqlite"
)

func TestGreetingDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	database, cleanup, err := db.Open(sqlite.Open(path))
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	repo := NewGreeting(database)
	if err := repo.Init(); err != nil {
		t.Fatal(err)
	}
	first, err := repo.Get()
	if err != nil || first.Prefix != "Hello" {
		t.Fatalf("greeting=%+v err=%v", first, err)
	}
	if err := database.Model(&first).Update("prefix", "Hi").Error; err != nil {
		t.Fatal(err)
	}
	second, secondCleanup, err := db.Open(sqlite.Open(path))
	if err != nil {
		t.Fatal(err)
	}
	defer secondCleanup()
	if err := NewGreeting(second).Init(); err != nil {
		t.Fatal(err)
	}
	got, err := NewGreeting(second).Get()
	if err != nil || got.Prefix != "Hi" {
		t.Fatalf("persisted=%+v err=%v", got, err)
	}
	var count int64
	database.Model(&model.Greeting{}).Count(&count)
	if count != 1 {
		t.Fatalf("count=%d", count)
	}
	if err := secondCleanup(); err != nil {
		t.Fatal(err)
	}
	if _, err := NewGreeting(second).Get(); err == nil {
		t.Fatal("expected closed database error")
	}
}
