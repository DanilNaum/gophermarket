package usecase

import (
	"context"

	"github.com/google/uuid"
)

func (u *usecase) GetOrders(ctx context.Context, userID uuid.UUID) ([]*Order, error) {
	orders, err := u.orderRepository.GetOrdersByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return ordersFromModel(orders), nil
}
