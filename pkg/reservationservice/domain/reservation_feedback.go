package domain

import (
	"time"
)

type ReservationFeedback struct {
	ID            int32     `db:"id"`
	ReservationID int32     `db:"reservation_id"`
	StatusID      int32     `db:"status_id"`
	SolutionID    int32     `db:"solution_id"`
	Rate          int32     `db:"rate"`
	Description   string    `db:"description"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
