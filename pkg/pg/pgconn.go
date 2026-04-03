package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type connection struct {
	masterPool *pgxpool.Pool
}

func (c *connection) Master() *pgxpool.Pool {
	return c.masterPool
}

func (c *connection) Close() {
	c.Master().Close()

}

func NewConnection(ctx context.Context, cnf *connConfig, log *zap.SugaredLogger) *connection {
	masterDsn := cnf.getDsn()

	masterPool := createPool(ctx, masterDsn, "master", log)
	if masterPool == nil {
		return nil
	}

	return &connection{
		masterPool: masterPool,
	}
}

func createPool(ctx context.Context, dsn, tp string, log *zap.SugaredLogger) *pgxpool.Pool {
	pg, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Errorf("сould not establish db %s connection %s", tp, err.Error())
		return nil
	}

	log.Info("msg", fmt.Sprintf("Database connection %s established", tp))
	return pg
}
