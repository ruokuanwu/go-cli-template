package app

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestInvalidConfig(t *testing.T) {
	t.Setenv("APP_LOG_LEVEL", "invalid")
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"greet"}, &stdout, &stderr, "dev"); code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "invalid log level") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
}

func TestInvalidDatabase(t *testing.T) {
	t.Setenv("APP_LOG_LEVEL", "info")
	t.Setenv("APP_DB_DSN", filepath.Join(t.TempDir(), "missing", "test.db"))
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"greet"}, &stdout, &stderr, "dev"); code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "database initialization failed") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
}

func TestRun(t *testing.T) {
	t.Setenv("APP_LOG_LEVEL", "info")
	t.Setenv("APP_DB_DSN", ":memory:")
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"greet"}, &stdout, &stderr, "dev"); code != 0 || stdout.String() != "Hello, world!\n" || stderr.Len() != 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
	stdout.Reset()
	if code := Run([]string{"missing"}, &stdout, &stderr, "dev"); code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
	var entry map[string]any
	if err := json.Unmarshal(stderr.Bytes(), &entry); err != nil {
		t.Fatal(err)
	}
	if entry["level"] != "error" || entry["msg"] != "command failed" || entry["error"] == nil {
		t.Fatalf("unexpected log: %v", entry)
	}
}
