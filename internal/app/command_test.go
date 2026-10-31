package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	t.Setenv("APP_LOG_LEVEL", "info")
	tests := []struct {
		name   string
		args   []string
		code   int
		stdout string
		stderr string
	}{
		{name: "default help", stdout: "Usage:"},
		{name: "help", args: []string{"--help"}, stdout: "Usage:"},
		{name: "short help", args: []string{"-h"}, stdout: "Usage:"},
		{name: "version", args: []string{"--version"}, stdout: "go-cli-template v1.2.3\n"},
		{name: "greet", args: []string{"greet"}, stdout: "Hello, world!\n"},
		{name: "named greet", args: []string{"greet", "--name", "Go"}, stdout: "Hello, Go!\n"},
		{name: "greet help", args: []string{"greet", "--help"}, stdout: "Usage: go-cli-template greet"},
		{name: "unknown command", args: []string{"missing"}, code: 2, stderr: "unknown command"},
		{name: "unknown flag", args: []string{"greet", "--missing"}, code: 2, stderr: "unknown flag"},
		{name: "missing flag value", args: []string{"greet", "--name"}, code: 2, stderr: "flag needs an argument"},
		{name: "empty name", args: []string{"greet", "--name", "  "}, code: 2, stderr: "name must not be empty"},
		{name: "extra argument", args: []string{"greet", "extra"}, code: 2, stderr: "unexpected arguments"},
		{name: "version extra argument", args: []string{"--version", "extra"}, code: 2, stderr: "unexpected arguments"},
		{name: "help extra argument", args: []string{"--help", "extra"}, code: 2, stderr: "unexpected arguments"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := Run(tt.args, &stdout, &stderr, "v1.2.3"); code != tt.code {
				t.Fatalf("exit code = %d, want %d; stderr: %s", code, tt.code, &stderr)
			}
			if tt.stdout != "" && !strings.Contains(stdout.String(), tt.stdout) {
				t.Errorf("stdout = %q, want %q", stdout.String(), tt.stdout)
			}
			if tt.stderr != "" && !strings.Contains(stderr.String(), tt.stderr) {
				t.Errorf("stderr = %q, want %q", stderr.String(), tt.stderr)
			}
			if tt.code == 0 && stderr.Len() != 0 {
				t.Errorf("unexpected stderr: %s", &stderr)
			}
			if tt.code != 0 && stdout.Len() != 0 {
				t.Errorf("unexpected stdout: %s", &stdout)
			}
		})
	}
}
