package config

import (
	"flag"

	"github.com/DanilNaum/gophermarket/internal/config/client"
	"github.com/DanilNaum/gophermarket/internal/config/db"
	"github.com/DanilNaum/gophermarket/internal/config/server"
	"go.uber.org/zap"
)

type dbConfig interface {
	GetDSN() string
	GetCryptoKey() string
}

type serverConf interface {
	HTTPServerHost() string
	ServerPort() (int, error)
}

type clientConfig interface {
	AccrualAddr() string
}

type config struct {
	dbConfig     dbConfig
	serverConf   serverConf
	clientConfig clientConfig
}

func NewConfig(log *zap.SugaredLogger) *config {
	dbConfigFlag := db.DBConfigFromFlags()
	serverConfigFlags := server.ServerConfigFromFlags()
	clientConfigFlags := client.ClientConfigFromFlags()

	flag.Parse()

	dbConfigEnv := db.DBConfigFromEnv(log)
	serverConfigEnv := server.ServerConfigFromEnv(log)
	clientConfigEnv := client.ClientConfigFromEnv(log)

	dbConfig := db.MergeDBConfigs(dbConfigEnv, dbConfigFlag, log)
	serverConfig := server.MergeServerConfigs(serverConfigEnv, serverConfigFlags, log)
	clientConfig := client.MergeClientConfigs(clientConfigEnv, clientConfigFlags, log)

	return &config{dbConfig: dbConfig, serverConf: serverConfig, clientConfig: clientConfig}
}

func (c *config) DBConfig() dbConfig {
	return c.dbConfig
}

func (c *config) ServerConfig() serverConf {
	return c.serverConf
}

func (c *config) ClientConfig() clientConfig {
	return c.clientConfig
}
