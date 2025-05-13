package usecase

import (
	"context"
	"errors"

	"github.com/DanilNaum/gophermarket/internal/repository/user"
	"github.com/DanilNaum/gophermarket/pkg/jwt"
)

func (u *usecase) Register(ctx context.Context, login, password string) (string, error) {
	// todo: hash password

	id, err := u.repo.CreateUser(ctx, &user.User{Login: login, PasswordHash: password, Salt: "empty_salt"})
	if err != nil {
		if errors.Is(err, user.ErrConflict) {
			return "", user.ErrConflict
		}
		return "", err
	}

	token, err := u.jwt.GenerateToken(jwt.User{ID: id})
	if err != nil {
		return "", err
	}

	return token, nil
}
