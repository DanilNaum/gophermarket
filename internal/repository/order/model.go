package order

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
