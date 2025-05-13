package db

import (
	"github.com/caarlos0/env/v6"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type dbConfig struct {
	DBDSN *string `env:"DATABASE_DSN"`
}

func DBConfigFromEnv(log logger) *dbConfig {
	c := &dbConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func (c *dbConfig) GetDSN() string {
	return *c.DBDSN
}
