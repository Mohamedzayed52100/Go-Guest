package domain

import (
	"time"
)

type ReservationTag struct {
	ID         int32                   `db:"id"`
	BranchID   int32                   `db:"branch_id"`
	Name       string                  `db:"name"`
	CategoryID int32                   `db:"category_id"`
	Category   *ReservationTagCategory `gorm:"foreignKey:CategoryID"`
	CreatedAt  time.Time               `db:"created_at"`
	UpdatedAt  time.Time               `db:"updated_at"`
}
