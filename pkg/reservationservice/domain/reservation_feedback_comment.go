package domain

import (
	"time"
)

type ReservationFeedbackComment struct {
	ID                    int32     `db:"id"`
	ReservationFeedbackID int32     `db:"reservation_feedback_id"`
	CreatorID             int32     `db:"creator_id"`
	Comment               string    `db:"comment"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}
