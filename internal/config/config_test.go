package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("APP_DB_DSN", "")
	if got := Load().DB.DSN; got != ":memory:" {
		t.Fatalf("default DSN = %q", got)
	}
	t.Setenv("APP_DB_DSN", "test.db")
	if got := Load().DB.DSN; got != "test.db" {
		t.Fatalf("DSN = %q", got)
	}
	t.Setenv("APP_LOG_LEVEL", "")
	cfg := Load()
	if cfg.Log.Level != "info" {
		t.Fatalf("default log level = %q", cfg.Log.Level)
	}
	t.Setenv("APP_LOG_LEVEL", "debug")
	if got := Load().Log.Level; got != "debug" {
		t.Fatalf("log level = %q", got)
	}
}
