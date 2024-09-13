package domain

import "time"

type ReservationStatus struct {
	ID        int       `db:"id"`
	BranchID  int32     `db:"branch_id"`
	Name      string    `db:"name"`
	Category  string    `db:"category"`
	Type      string    `db:"type"`
	Color     string    `db:"color"`
	Icon      string    `db:"icon"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
