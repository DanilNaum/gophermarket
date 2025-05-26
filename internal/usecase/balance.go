package usecase

import (
	"context"

	"github.com/google/uuid"
)

func (u *usecase) Balance(ctx context.Context, userID uuid.UUID) (current int, withdrawn int, err error) {
	orderSum, withdrawalSum, err := u.withdrawalRepository.GetBalance(ctx, userID)
	if err != nil {
		return 0, 0, err
	}

	return orderSum - withdrawalSum, withdrawalSum, nil
}
