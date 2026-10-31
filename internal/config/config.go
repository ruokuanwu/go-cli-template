package config

import "os"

type Config struct {
	Log Log
	DB  Database
}

type Database struct {
	DSN string
}

type Log struct {
	Level string
}

func Load() Config {
	level := os.Getenv("APP_LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	dsn := os.Getenv("APP_DB_DSN")
	if dsn == "" {
		dsn = ":memory:"
	}
	return Config{Log: Log{Level: level}, DB: Database{DSN: dsn}}
}
