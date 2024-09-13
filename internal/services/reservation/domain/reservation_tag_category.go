package domain

import "time"

type ReservationTagCategory struct {
	ID             int32             `db:"id"`
	BranchID       int32             `db:"branch_id"`
	Name           string            `db:"name"`
	Color          string            `db:"color"`
	Classification string            `db:"classification"`
	OrderIndex     int32             `db:"order_index"`
	Tags           []*ReservationTag `gorm:"-"`
	IsDisabled     bool              `db:"is_disabled"`
	CreatedAt      time.Time         `db:"created_at"`
	UpdatedAt      time.Time         `db:"updated_at"`
}
