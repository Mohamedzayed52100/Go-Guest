package domain

import "time"

type SimpleReservationFeedback struct {
	ID          int32     `db:"id"`
	Rate        int32     `db:"rate"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (SimpleReservationFeedback) TableName() string {
	return "reservation_feedbacks"
}
