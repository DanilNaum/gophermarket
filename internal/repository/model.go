package repository

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        string
	UserID    uuid.UUID
	Status    string
	Accrual   *int
	CreatedAt time.Time
}

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	Salt         string
}

type Withdrawal struct {
	ID        string
	UserID    uuid.UUID
	Accrual   *int
	CreatedAt time.Time
}
