package usecase

import (
	"time"

	"github.com/DanilNaum/gophermarket/internal/repository/order"
	"github.com/DanilNaum/gophermarket/internal/repository/withdrawal"
)

type Order struct {
	ID         string
	Status     string
	Accrual    *int
	UploadedAt time.Time
}

func orderFromModel(order *order.Order) *Order {
	return &Order{
		ID:         order.ID,
		Status:     order.Status,
		Accrual:    order.Accrual,
		UploadedAt: order.CreatedAt,
	}
}

func ordersFromModel(orders []*order.Order) []*Order {
	ordersModel := make([]*Order, 0, len(orders))
	for _, order := range orders {
		ordersModel = append(ordersModel, orderFromModel(order))
	}
	return ordersModel
}

type Withdraw struct {
	Order       string
	Sum         *int
	ProcessedAt time.Time
}

func withdrawalFromModel(w *withdrawal.Withdrawal) *Withdraw {
	return &Withdraw{
		Order:       w.ID,
		Sum:         w.Accrual,
		ProcessedAt: w.CreatedAt,
	}
}
func withdrawalsFromModel(ws []*withdrawal.Withdrawal) []*Withdraw {
	wsModel := make([]*Withdraw, 0, len(ws))
	for _, w := range ws {
		wsModel = append(wsModel, withdrawalFromModel(w))
	}
	return wsModel
}
