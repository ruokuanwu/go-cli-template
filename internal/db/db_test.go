package db

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

type testDialector struct {
	gorm.Dialector
	err error
}

func (testDialector) Name() string                { return "test" }
func (d testDialector) Initialize(*gorm.DB) error { return d.err }

func TestOpen(t *testing.T) {
	t.Run("initialization error", func(t *testing.T) {
		want := errors.New("connection failed")
		_, err := Open(testDialector{err: want})
		if !errors.Is(err, want) {
			t.Fatalf("error = %v, want wrapped %v", err, want)
		}
	})
	t.Run("success", func(t *testing.T) {
		database, err := Open(testDialector{})
		if err != nil {
			t.Fatal(err)
		}
		if database == nil {
			t.Fatal("database is nil")
		}
	})
}
