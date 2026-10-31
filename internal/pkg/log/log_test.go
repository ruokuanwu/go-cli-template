package log

import (
	"bytes"
	"strings"
	"testing"

	"example.com/go-cli-template/internal/config"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	previous := zap.L()
	defer zap.ReplaceGlobals(previous)
	var output bytes.Buffer
	if err := Init(config.Log{Level: "error"}, &output); err != nil {
		t.Fatal(err)
	}
	zap.L().Info("hidden")
	zap.L().Error("visible")
	if strings.Contains(output.String(), "hidden") || !strings.Contains(output.String(), "visible") {
		t.Fatalf("unexpected logs: %s", &output)
	}
	logger := zap.L()
	if err := Init(config.Log{Level: "invalid"}, &output); err == nil {
		t.Fatal("expected invalid level error")
	}
	if zap.L() != logger {
		t.Fatal("invalid config replaced global logger")
	}
}
