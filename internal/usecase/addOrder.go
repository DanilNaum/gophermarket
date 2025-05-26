package usecase

import (
	"context"
	"errors"

	"github.com/DanilNaum/gophermarket/internal/repository"
	"github.com/DanilNaum/gophermarket/internal/usecase/model"
	"github.com/google/uuid"
)

func (u *usecase) AddOrder(ctx context.Context, userID uuid.UUID, orderID string) error {
	err := u.orderRepository.CreateOrder(ctx, &repository.Order{
		ID:      orderID,
		UserID:  userID,
		Status:  model.StatusNEW,
		Accrual: nil,
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrExists):
			return ErrExists
		case errors.Is(err, repository.ErrConflict):
			return ErrConflict
		default:
			return ErrUnexpected
		}
	}

	u.orderService.NewOrder(orderID)

	return nil
}
