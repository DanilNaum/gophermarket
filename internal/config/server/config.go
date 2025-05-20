package server

import (
	"flag"
	"net"
	"strconv"

	"net/url"

	"github.com/caarlos0/env/v6"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type serverConfig struct {
	Host *string `env:"RUN_ADDRESS"`
}

func ServerConfigFromFlags() *serverConfig {
	host := flag.String("a", "localhost:8080", "host:port")

	return &serverConfig{
		Host: host,
	}
}

func ServerConfigFromEnv(log logger) *serverConfig {
	c := &serverConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func MergeServerConfigs(envConfig, flagsConfig *serverConfig, log logger) *serverConfig {
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

func (c *serverConfig) ValidateServerConfig(log logger) {

	_, err := url.Parse(*c.Host)
	if err != nil {
		log.Fatalf("invalid host: %s", *c.Host)
		return
	}

}

func (c *serverConfig) HTTPServerHost() string {
	return *c.Host
}

func (c *serverConfig) ServerPort() (int, error) {
	_, port, err := net.SplitHostPort(*c.Host)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(port)
}
