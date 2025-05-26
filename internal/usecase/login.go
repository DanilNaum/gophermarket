package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/DanilNaum/gophermarket/internal/repository"
	"github.com/DanilNaum/gophermarket/pkg/jwt"
)

func (u *usecase) Login(ctx context.Context, login, password string) (string, error) {
	usr, err := u.userRepository.GetUser(ctx, login)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%w: %w", ErrUnexpected, err)

	}

	password, err = u.crypto.Encode(password)
	if err != nil {
		return "", ErrUnexpected
	}

	if usr.PasswordHash != password {
		return "", ErrInvalidPassword
	}

	token, err := u.jwt.GenerateToken(jwt.User{ID: usr.ID})
	if err != nil {
		return "", err
	}

	return token, nil
}
