package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
)

func NewWithdrawalStorage(conn connection) *storage {
	return &storage{conn: conn}
}

func (s *storage) AddWithdrawal(ctx context.Context, w *Withdrawal) error {
	query := `WITH balances AS (
   	 		SELECT
        		(SELECT COALESCE(SUM(accrual), 0) FROM orders WHERE user_id = $1) AS total_orders,
        		(SELECT COALESCE(SUM(accrual), 0) FROM withdrawal WHERE user_id = $1) AS total_withdrawals
			)
			INSERT INTO withdrawal (id, user_id, accrual)
			SELECT $2, $1, $3
			FROM balances
			WHERE (total_orders - total_withdrawals) >= $3`
	cTag, err := s.conn.Master().Exec(ctx, query, w.UserID, w.ID, w.Accrual)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return ErrInvalidUser
			}
		}
		return err
	}
	if cTag.RowsAffected() == 0 {
		return ErrNotEnough
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

func (s *storage) GetBalance(ctx context.Context, userID uuid.UUID) (int, int, error) {
	var (
		total_withdrawals int
		total_orders      int
	)
	query := `SELECT
        		(SELECT COALESCE(SUM(accrual), 0) FROM orders WHERE user_id = $1) AS total_orders,
        		(SELECT COALESCE(SUM(accrual), 0) FROM withdrawal WHERE user_id = $1) AS total_withdrawals`
	err := s.conn.Master().QueryRow(ctx, query, userID).Scan(&total_orders, &total_withdrawals)
	if err != nil {
		return 0, 0, err
	}
	return total_orders, total_withdrawals, nil
}
