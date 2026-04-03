package usecase

import (
	"context"

	"github.com/DanilNaum/gophermarket/internal/usecase/model"
	"github.com/google/uuid"
)

func (u *usecase) GetOrders(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	orders, err := u.orderRepository.GetOrdersByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return model.OrdersFromModel(orders), nil
}
