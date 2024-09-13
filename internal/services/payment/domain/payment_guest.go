package domain

type PaymentGuest struct {
	ID          int32
	FirstName   string
	LastName    string
	PhoneNumber string
	Address     string
	Email       *string
}
