package domain

import "time"

type ReservationVisitor struct {
	ID            int32     `db:"id"`
	ReservationID int32     `db:"reservation_id"`
	GuestID       int32     `db:"guest_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
