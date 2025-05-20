package usecase

import (
	"context"

	"github.com/DanilNaum/gophermarket/internal/repository/order"
	"github.com/DanilNaum/gophermarket/internal/repository/user"
	"github.com/DanilNaum/gophermarket/internal/repository/withdrawal"
	"github.com/DanilNaum/gophermarket/pkg/jwt"
	"github.com/google/uuid"
)

type userRepository interface {
	CreateUser(ctx context.Context, user *user.User) (uuid.UUID, error)
	GetUser(ctx context.Context, login string) (*user.User, error)
}
type orderRepository interface {
	CreateOrder(ctx context.Context, order *order.Order) error
	GetOrdersByUserId(ctx context.Context, userID uuid.UUID) ([]*order.Order, error)
	GetUserOrdersSum(ctx context.Context, userID uuid.UUID) (int, error)
}

type withdrawalRepository interface {
	AddWithdrawal(ctx context.Context, w *withdrawal.Withdrawal) error
	GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]*withdrawal.Withdrawal, error)
	GetUserWithdrawalsSum(ctx context.Context, userID uuid.UUID) (int, error)
}

type jwtManager interface {
	GenerateToken(user jwt.User) (string, error)
	ParseToken(tokenString string) (*jwt.User, error)
}
type orderService interface {
	NewOrder(orderID string)
}

type usecase struct {
	userRepository       userRepository
	orderRepository      orderRepository
	withdrawalRepository withdrawalRepository
	orderService         orderService
	jwt                  jwtManager
}

func NewUsecase(userRepository userRepository, orderRepository orderRepository, jwt jwtManager, orderService orderService, withdrawalRepository withdrawalRepository) (*usecase, error) {

	return &usecase{
		userRepository:       userRepository,
		orderRepository:      orderRepository,
		withdrawalRepository: withdrawalRepository,
		jwt:                  jwt,
		orderService:         orderService,
	}, nil
}
