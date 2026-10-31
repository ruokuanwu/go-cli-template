package config

import "testing"

func TestLoad(t *testing.T) {
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
