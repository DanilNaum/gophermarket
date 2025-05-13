package usecase

import (
	"context"

	"github.com/DanilNaum/gophermarket/internal/repository/user"
	"github.com/DanilNaum/gophermarket/pkg/jwt"
	"github.com/google/uuid"
)

type repository interface {
	CreateUser(ctx context.Context, user *user.User) (uuid.UUID, error)
	GetUser(ctx context.Context, login string) (*user.User, error)
}

type jwtManager interface {
	GenerateToken(user jwt.User) (string, error)
	ParseToken(tokenString string) (*jwt.User, error)
}

type usecase struct {
	repo repository
	jwt  jwtManager
}

func NewUsecase(repo repository, jwt jwtManager) (*usecase, error) {
	// if repo == nil {
	// 	return nil, errors.New("repository is nil")
	// }
	return &usecase{repo: repo, jwt: jwt}, nil
}
