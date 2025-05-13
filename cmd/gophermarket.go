package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/DanilNaum/gophermarket/internal/api/restapi"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations"
	"github.com/DanilNaum/gophermarket/internal/config"
	"github.com/DanilNaum/gophermarket/internal/repository/user"
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

	userStorage := user.NewUserStorage(pgConn)

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		panic(err)
	}

	api := operations.NewGopherMarketLoyaltySystemAPIAPI(swaggerSpec)
	api.UseSwaggerUI()

	jwt := jwt.NewJWTManager(jwt.WithTokenExpiration(time.Minute), jwt.WithSecretKey([]byte("secret")))

	usecase, err := usecase.NewUsecase(userStorage, jwt)
	if err != nil {
		return err
	}

	configurator, err := transport.NewConfigurator(usecase)
	if err != nil {
		return err
	}
	configurator.Configure(api)

	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = _serverPort

	err = server.Serve()
	if err != nil {
		return err
	}
	return nil

}
