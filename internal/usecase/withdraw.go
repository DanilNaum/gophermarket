package usecase

import (
	"context"
	"errors"

	"github.com/DanilNaum/gophermarket/internal/repository"
	"github.com/DanilNaum/gophermarket/internal/usecase/model"
	"github.com/google/uuid"
)

func (u *usecase) Withdraw(ctx context.Context, userID uuid.UUID, orderID string, accrual int) error {

	err := u.withdrawalRepository.AddWithdrawal(ctx, &repository.Withdrawal{
		ID:      orderID,
		UserID:  userID,
		Accrual: &accrual,
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotEnough):
			return ErrNotEnough
		default:
			return err
		}

	}

	return nil
}

func (u *usecase) GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]*model.Withdrawal, error) {
	ws, err := u.withdrawalRepository.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.WithdrawalsFromModel(ws), nil
}
