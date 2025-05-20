package usecase

import (
	"context"

	"github.com/DanilNaum/gophermarket/internal/repository/withdrawal"
	"github.com/google/uuid"
)

func (u *usecase) Withdraw(ctx context.Context, userID uuid.UUID, orderID string, accrual int) error {
	orderSum, err := u.orderRepository.GetUserOrdersSum(ctx, userID)
	if err != nil {
		return err
	}
	withdrawalSum, err := u.withdrawalRepository.GetUserWithdrawalsSum(ctx, userID)
	if err != nil {
		return err
	}

	if orderSum-withdrawalSum < accrual {
		return ErrNotEnough
	}

	err = u.withdrawalRepository.AddWithdrawal(ctx, &withdrawal.Withdrawal{
		ID:      orderID,
		UserID:  userID,
		Accrual: &accrual,
	})
	if err != nil {
		return err
	}

	return nil
}
