package usecase

import (
	"context"

	repository "github.com/DanilNaum/gophermarket/internal/repository"
	"github.com/DanilNaum/gophermarket/pkg/jwt"
	"github.com/google/uuid"
)

type userRepository interface {
	CreateUser(ctx context.Context, user *repository.User) (uuid.UUID, error)
	GetUser(ctx context.Context, login string) (*repository.User, error)
}
type orderRepository interface {
	CreateOrder(ctx context.Context, order *repository.Order) error
	GetOrdersByUserId(ctx context.Context, userID uuid.UUID) ([]*repository.Order, error)
	GetUserOrdersSum(ctx context.Context, userID uuid.UUID) (int, error)
}

type withdrawalRepository interface {
	AddWithdrawal(ctx context.Context, w *repository.Withdrawal) error
	GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]*repository.Withdrawal, error)
	GetBalance(ctx context.Context, userID uuid.UUID) (int, int, error)
}

type jwtManager interface {
	GenerateToken(user jwt.User) (string, error)
	ParseToken(tokenString string) (*jwt.User, error)
}
type orderService interface {
	NewOrder(orderID string)
}
type crypto interface {
	Encode(password string) (string, error)
}

type usecase struct {
	userRepository       userRepository
	orderRepository      orderRepository
	withdrawalRepository withdrawalRepository
	orderService         orderService
	jwt                  jwtManager
	crypto               crypto
}

func NewUsecase(userRepository userRepository, orderRepository orderRepository, jwt jwtManager, orderService orderService, withdrawalRepository withdrawalRepository, crypto crypto) (*usecase, error) {

	return &usecase{
		userRepository:       userRepository,
		orderRepository:      orderRepository,
		withdrawalRepository: withdrawalRepository,
		jwt:                  jwt,
		orderService:         orderService,
		crypto:               crypto,
	}, nil
}
