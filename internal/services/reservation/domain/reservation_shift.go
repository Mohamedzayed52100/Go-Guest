package domain

type ReservationShift struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

func (ReservationShift) TableName() string {
	return "shifts"
}
