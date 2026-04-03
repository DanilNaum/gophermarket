package db

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type dbConfig struct {
	DBDSN     *string `env:"DATABASE_URI"`
	CryptoKey *string `env:"CRYPTO_KEY" envDefault:"123456"`
}

func DBConfigFromFlags() *dbConfig {
	dsn := flag.String("d", "", "dsn")

	return &dbConfig{
		DBDSN: dsn,
	}
}

func DBConfigFromEnv(log *zap.SugaredLogger) *dbConfig {
	c := &dbConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func MergeDBConfigs(envConfig, flagsConfig *dbConfig, log *zap.SugaredLogger) *dbConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if envConfig.DBDSN == nil {
		envConfig.DBDSN = flagsConfig.DBDSN
	}

	return envConfig
}

func (c *dbConfig) GetDSN() string {
	return *c.DBDSN
}

func (c *dbConfig) GetCryptoKey() string {
	return *c.CryptoKey
}
