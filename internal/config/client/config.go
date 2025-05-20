package client

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type clientConfig struct {
	Host *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func ClientConfigFromFlags() *clientConfig {
	host := flag.String("r", "localhost:8080", "host:port")

	return &clientConfig{
		Host: host,
	}
}

func ClientConfigFromEnv(log logger) *clientConfig {
	c := &clientConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func MergeClientConfigs(envConfig, flagsConfig *clientConfig, log logger) *clientConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if envConfig.Host == nil {
		envConfig.Host = flagsConfig.Host
	}

	return envConfig
}

func (c *clientConfig) AccrualAddr() string {
	return *c.Host
}
