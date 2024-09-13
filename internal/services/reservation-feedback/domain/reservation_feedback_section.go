package domain

import "time"

type ReservationFeedbackSection struct {
	ID        int32     `db:"id"`
	Name      string    `db:"name"`
	BranchID  int32     `db:"branch_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
