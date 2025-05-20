package usecase

import (
	"context"
	"errors"

	"github.com/DanilNaum/gophermarket/internal/repository/order"
	"github.com/google/uuid"
)

func (u *usecase) AddOrder(ctx context.Context, userID uuid.UUID, orderID string) error {
	err := u.orderRepository.CreateOrder(ctx, &order.Order{
		ID:      orderID,
		UserID:  userID,
		Status:  "NEW",
		Accrual: nil,
	})
	if err != nil {
		switch {
		case errors.Is(err, order.ErrExists):
			return ErrExists
		case errors.Is(err, order.ErrConflict):
			return ErrConflict
		default:
			return ErrUnexpected
		}
	}

	u.orderService.NewOrder(orderID)

	return nil
}
