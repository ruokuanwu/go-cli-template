package config

import "os"

type Config struct {
	Log Log
}

type Log struct {
	Level string
}

func Load() Config {
	level := os.Getenv("APP_LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	return Config{Log: Log{Level: level}}
}
