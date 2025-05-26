package main

import (
	"context"
	"errors"
	"github.com/DanilNaum/gophermarket/pkg/crypto"
	"os"
	"time"

	"github.com/DanilNaum/gophermarket/internal/api/restapi"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations"
	"github.com/DanilNaum/gophermarket/internal/config"
	"github.com/DanilNaum/gophermarket/internal/repository"
	orderService "github.com/DanilNaum/gophermarket/internal/service/order"
	"github.com/DanilNaum/gophermarket/internal/transport"
	"github.com/DanilNaum/gophermarket/internal/usecase"
	"github.com/DanilNaum/gophermarket/pkg/jwt"
	"github.com/DanilNaum/gophermarket/pkg/migration"
	"github.com/DanilNaum/gophermarket/pkg/pg"
	"github.com/go-openapi/loads"
	"go.uber.org/zap"
)

const _serverPort = 80

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}

	defer logger.Sync()

	sugarLogger := logger.Sugar()

	sugarLogger.Info("App is running...")

	err = run(sugarLogger)
	if err != nil {
		sugarLogger.Error(err)
		os.Exit(1)
	}
	os.Exit(0)

}

func run(log *zap.SugaredLogger) error {
	ctx := context.Background()

	conf := config.NewConfig(log)

	migrator := migration.NewMigrator(conf.DBConfig().GetDSN(), migration.WithRelativePath("migrations"))
	err := migrator.Migrate()
	if err != nil {
		return err
	}

	pgConf := pg.NewConnConfigFromDsnString(conf.DBConfig().GetDSN())

	pgConn := pg.NewConnection(ctx, pgConf, log)
	if pgConn == nil {
		return errors.New("pg connection is nil")
	}
	defer pgConn.Close()

	storage := repository.NewStorage(pgConn)

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		panic(err)
	}

	api := operations.NewGopherMarketLoyaltySystemAPIAPI(swaggerSpec)
	api.UseSwaggerUI()

	jwt := jwt.NewJWTManager(jwt.WithTokenExpiration(time.Minute), jwt.WithSecretKey([]byte("secret")))

	orderService := orderService.NewOrderService(storage, 10, conf.ClientConfig().AccrualAddr())
	go orderService.Start(ctx)

	crypto := crypto.NewCrypto(conf.DBConfig().GetCryptoKey())

	usecase, err := usecase.NewUsecase(storage, storage, jwt, orderService, storage, crypto)
	if err != nil {
		return err
	}

	configurator, err := transport.NewConfigurator(usecase, log)
	if err != nil {
		return err
	}
	configurator.Configure(api)

	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = _serverPort
	if port, err := conf.ServerConfig().ServerPort(); err == nil {
		server.Port = port
	}

	err = server.Serve()
	if err != nil {
		return err
	}
	return nil

}
