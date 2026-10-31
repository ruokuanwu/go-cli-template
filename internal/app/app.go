package app

import (
	"fmt"
	"io"

	"example.com/go-cli-template/internal/config"
	"example.com/go-cli-template/internal/db"
	"example.com/go-cli-template/internal/pkg/log"
	"example.com/go-cli-template/internal/repository"
	"example.com/go-cli-template/internal/service"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
)

func Run(args []string, stdout, stderr io.Writer, version string) int {
	cfg := config.Load()
	if err := log.Init(cfg.Log, stderr); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	defer zap.L().Sync()

	database, cleanup, err := db.Open(sqlite.Open(cfg.DB.DSN))
	if err != nil {
		zap.L().Error("database initialization failed", zap.Error(err))
		return 2
	}
	defer cleanup()

	repo := repository.NewGreeting(database)
	if err := repo.Init(); err != nil {
		zap.L().Error("repository initialization failed", zap.Error(err))
		return 2
	}
	greetings := service.NewGreeting(repo)
	if err := execute(args, stdout, stderr, version, greetings); err != nil {
		zap.L().Error("command failed", zap.Error(err), zap.Int("exit_code", 2))
		return 2
	}
	return 0
}
