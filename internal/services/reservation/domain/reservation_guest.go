package domain

type ReservationGuest struct {
	ID          int32  `db:"id"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	PhoneNumber string `db:"phone_number"`
}

func (ReservationGuest) TableName() string {
	return "guests"
}
