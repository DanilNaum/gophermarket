package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
)

func NewOrderStorage(conn connection) *storage {
	return &storage{conn: conn}
}

func (s *storage) CreateOrder(ctx context.Context, order *Order) error {
	var existingUserID uuid.UUID

	query := `
		INSERT INTO orders (id, user_id, status, accrual)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET status = EXCLUDED.status
		RETURNING user_id
	`

	err := s.conn.Master().QueryRow(ctx, query,
		order.ID, order.UserID, order.Status, order.Accrual,
	).Scan(&existingUserID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return ErrInvalidUser
		}
		return err
	}

	if existingUserID != order.UserID {
		return ErrConflict
	}

	return nil
}

const expectedOrdersCount = 10

func (s *storage) GetOrdersByUserId(ctx context.Context, userID uuid.UUID) ([]*Order, error) {
	orders := make([]*Order, 0, expectedOrdersCount)
	query := `SELECT id, status, accrual, created_at FROM orders WHERE user_id = $1`
	rows, err := s.conn.Master().Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		order := &Order{}
		err = rows.Scan(&order.ID, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil

}

func (s *storage) UpdateOrder(ctx context.Context, order *Order) error {
	query := `UPDATE orders SET status = $1, accrual = $2 WHERE id = $3`
	_, err := s.conn.Master().Exec(ctx, query, order.Status, order.Accrual, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetUserOrdersSum(ctx context.Context, userID uuid.UUID) (int, error) {
	var su int
	query := `SELECT COALESCE(SUM(accrual), 0) FROM orders WHERE user_id = $1`
	err := s.conn.Master().QueryRow(ctx, query, userID).Scan(&su)
	if err != nil {
		return 0, err
	}
	return su, nil
}
