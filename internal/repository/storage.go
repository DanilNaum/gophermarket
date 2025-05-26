package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type connection interface {
	Master() *pgxpool.Pool
	Close()
}
type storage struct {
	conn connection
}

func NewStorage(conn connection) *storage {
	return &storage{conn: conn}
}
