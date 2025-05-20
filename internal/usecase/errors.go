package usecase

import "errors"

var (
	ErrConflict        = errors.New("conflict")
	ErrNotFound        = errors.New("not found")
	ErrUnexpected      = errors.New("unexpected")
	ErrInvalidPassword = errors.New("invalid password")
	ErrExists          = errors.New("exists")
	ErrNotEnough       = errors.New("not enough")
)
