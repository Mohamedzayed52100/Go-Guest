package domain

import "time"

type SpecialOccasion struct {
	ID        int       `db:"id"`
	BranchID  int       `db:"branch_id"`
	Name      string    `db:"name"`
	Color     string    `db:"color"`
	Icon      string    `db:"icon"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
