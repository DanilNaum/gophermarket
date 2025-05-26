package repository

import "errors"

var (
	ErrConflict    = errors.New("conflict")
	ErrInvalidUser = errors.New("not found")
	ErrExists      = errors.New("exists")
)

var (
	ErrNotFound  = errors.New("user not found")
	ErrNotEnough = errors.New("no enough")
)
