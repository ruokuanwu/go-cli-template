package app

import (
	"fmt"
	"io"

	"example.com/go-cli-template/internal/config"
	"example.com/go-cli-template/internal/pkg/log"
	"example.com/go-cli-template/internal/repository"
	"example.com/go-cli-template/internal/service"
	"go.uber.org/zap"
)

func Run(args []string, stdout, stderr io.Writer, version string) int {
	cfg := config.Load()
	if err := log.Init(cfg.Log, stderr); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	defer zap.L().Sync()

	greetings := service.NewGreeting(repository.Greeting{})
	if err := execute(args, stdout, stderr, version, greetings); err != nil {
		zap.L().Error("command failed", zap.Error(err), zap.Int("exit_code", 2))
		return 2
	}
	return 0
}
