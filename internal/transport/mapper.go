package transport

import (
	"github.com/DanilNaum/gophermarket/internal/api/models"
	uc "github.com/DanilNaum/gophermarket/internal/usecase/model"
	"github.com/go-openapi/strfmt"
)

func orderFromUCModel(order *uc.Order) *models.Order {
	o := &models.Order{
		Number:     order.ID,
		Status:     order.Status,
		UploadedAt: strfmt.DateTime(order.UploadedAt),
	}
	if order.Accrual != nil {
		o.Accrual = float64(*order.Accrual)
	}
	return o
}

func ordersFromUCModel(orders []*uc.Order) []*models.Order {
	ordersModel := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		ordersModel = append(ordersModel, orderFromUCModel(order))
	}
	return ordersModel
}

func withdrawalFromUCModel(withdraw *uc.Withdrawal) *models.Withdrawal {
	w := &models.Withdrawal{
		Order:       withdraw.Order,
		ProcessedAt: strfmt.DateTime(withdraw.ProcessedAt),
	}
	if withdraw.Sum != nil {
		w.Sum = float64(*withdraw.Sum)
	}
	return w
}
func withdrawsFromUCModel(withdraws []*uc.Withdrawal) []*models.Withdrawal {
	ws := make([]*models.Withdrawal, 0, len(withdraws))
	for _, w := range withdraws {
		ws = append(ws, withdrawalFromUCModel(w))
	}
	return ws
}
