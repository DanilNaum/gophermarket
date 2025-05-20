package withdrawal

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
)

type connection interface {
	Master() *pgxpool.Pool
	Close()
}
type storage struct {
	conn connection
}

func NewWithdrawalStorage(conn connection) *storage {
	return &storage{conn: conn}
}

func (s *storage) AddWithdrawal(ctx context.Context, w *Withdrawal) error {
	query := "INSERT INTO withdrawal (order_id, user_id, accrual) VALUES ($1, $2, $3)"
	_, err := s.conn.Master().Exec(ctx, query, w.ID, w.UserID, w.Accrual)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return ErrInvalidUser
			}
		}
		return err
	}
	return nil
}

const expectedWithdrawalsNum = 10

func (s *storage) GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]*Withdrawal, error) {
	ws := make([]*Withdrawal, 0, expectedWithdrawalsNum)

	query := `SELECT order_id, accrual, processed_at FROM withdrawal WHERE user_id = $1`
	rows, err := s.conn.Master().Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		withdrawal := &Withdrawal{}
		err := rows.Scan(&withdrawal.ID, &withdrawal.Accrual, &withdrawal.CreatedAt)
		if err != nil {
			return nil, err
		}
		ws = append(ws, withdrawal)
	}
	return ws, nil
}

func (s *storage) GetUserWithdrawalsSum(ctx context.Context, userID uuid.UUID) (int, error) {
	var su int
	query := `SELECT COALESCE(SUM(accrual), 0) FROM withdrawal WHERE user_id = $1`
	err := s.conn.Master().QueryRow(ctx, query, userID).Scan(&su)
	if err != nil {
		return 0, err
	}
	return su, nil
}
