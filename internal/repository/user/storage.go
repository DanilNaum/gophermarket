package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type connection interface {
	Master() *pgxpool.Pool
	Close()
}
type storage struct {
	conn connection
}

func NewUserStorage(conn connection) *storage {
	return &storage{conn: conn}
}

func (s *storage) CreateUser(ctx context.Context, user *User) (uuid.UUID, error) {
	query := `INSERT INTO users (login, password_hash, password_salt) VALUES ($1, $2, $3) RETURNING id `

	err := s.conn.Master().QueryRow(ctx, query, user.Login, user.PasswordHash, user.Salt).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return uuid.UUID{}, ErrConflict
			}
			return uuid.UUID{}, err
		}
	}

	return user.ID, nil
}

func (s *storage) GetUser(ctx context.Context, login string) (*User, error) {
	user := &User{}

	query := `SELECT id, login, password_hash, password_salt FROM users WHERE login = $1`

	err := s.conn.Master().QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Salt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
