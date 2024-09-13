package domain

import "time"

type ReservationTagsAssignment struct {
	ID            int32     `db:"id"`
	TagID         int32     `db:"tag_id"`
	ReservationID int32     `db:"reservation_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
