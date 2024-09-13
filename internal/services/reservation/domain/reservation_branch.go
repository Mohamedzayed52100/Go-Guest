package domain

type ReservationBranch struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (ReservationBranch) TableName() string {
	return "branches"
}