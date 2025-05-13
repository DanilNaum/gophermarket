package usecase

import "github.com/google/uuid"

func (u *usecase) ParseToken(token string) (uuid.UUID, error) {
	user, err := u.jwt.ParseToken(token)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}
