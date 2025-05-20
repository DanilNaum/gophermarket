package usecase

import (
	"context"

	"github.com/google/uuid"
)

func (u *usecase) WithdrawHistory(ctx context.Context, userID uuid.UUID) ([]*Withdraw, error) {
	ws, err := u.withdrawalRepository.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	return withdrawalsFromModel(ws), nil
}
