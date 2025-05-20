package order

import "errors"

var (
	ErrConflict    = errors.New("conflict")
	ErrInvalidUser = errors.New("not found")
	ErrExists      = errors.New("exists")
)
