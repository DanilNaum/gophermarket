package model

import (
	"time"

	"github.com/DanilNaum/gophermarket/internal/repository"
)

const (
	StatusNEW = "NEW"
)

type Order struct {
	ID         string
	Status     string
	Accrual    *int
	UploadedAt time.Time
}

func orderFromModel(order *repository.Order) *Order {
	return &Order{
		ID:         order.ID,
		Status:     order.Status,
		Accrual:    order.Accrual,
		UploadedAt: order.CreatedAt,
	}
}

func OrdersFromModel(orders []*repository.Order) []*Order {
	ordersModel := make([]*Order, 0, len(orders))
	for _, order := range orders {
		ordersModel = append(ordersModel, orderFromModel(order))
	}
	return ordersModel
}

type Withdrawal struct {
	Order       string
	Sum         *int
	ProcessedAt time.Time
}

func withdrawalFromModel(w *repository.Withdrawal) *Withdrawal {
	return &Withdrawal{
		Order:       w.ID,
		Sum:         w.Accrual,
		ProcessedAt: w.CreatedAt,
	}
}
func WithdrawalsFromModel(ws []*repository.Withdrawal) []*Withdrawal {
	wsModel := make([]*Withdrawal, 0, len(ws))
	for _, w := range ws {
		wsModel = append(wsModel, withdrawalFromModel(w))
	}
	return wsModel
}
