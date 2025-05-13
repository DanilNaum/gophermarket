package user

import "errors"

var (
	ErrConflict = errors.New("user already exists")
	ErrNotFound = errors.New("user not found")
)
