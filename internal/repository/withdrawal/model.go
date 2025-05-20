package withdrawal

import (
	"time"

	"github.com/google/uuid"
)

type Withdrawal struct {
	ID        string
	UserID    uuid.UUID
	Accrual   *int
	CreatedAt time.Time
}
