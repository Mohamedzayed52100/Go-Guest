package domain

import "time"

type ReservationFeedbackSectionAssignment struct {
	ID         int32     `db:"id"`
	FeedbackID int32     `db:"feedback_id"`
	SectionID  int32     `db:"section_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
