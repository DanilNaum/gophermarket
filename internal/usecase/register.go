package usecase

import (
	"context"
	"errors"

	"github.com/DanilNaum/gophermarket/internal/repository"
	"github.com/DanilNaum/gophermarket/pkg/jwt"
)

func (u *usecase) Register(ctx context.Context, login, password string) (string, error) {
	password, err := u.crypto.Encode(password)
	if err != nil {
		return "", ErrUnexpected
	}

	id, err := u.userRepository.CreateUser(ctx, &repository.User{Login: login, PasswordHash: password, Salt: "empty_salt"})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return "", repository.ErrConflict
		}
		return "", err
	}

	token, err := u.jwt.GenerateToken(jwt.User{ID: id})
	if err != nil {
		return "", err
	}

	return token, nil
}
