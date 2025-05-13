package config

import "github.com/DanilNaum/gophermarket/internal/config/db"

type logger interface {
	Fatalf(format string, v ...any)
}

type dbConfig interface {
	GetDSN() string
}

type config struct {
	dbConfig dbConfig
}

func NewConfig(log logger) *config {
	dbConfig := db.DBConfigFromEnv(log)
	return &config{dbConfig: dbConfig}
}

func (c *config) DBConfig() dbConfig {
	return c.dbConfig
}
