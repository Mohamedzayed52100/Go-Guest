package domain

type Invoice struct {
	ID               int32   `db:"id"`
	InvoiceID        string  `db:"invoice_id"`
	PaymentRequestID int32   `db:"payment_request_id"`
	CustomerID       string  `db:"customer_id"`
	Status           string  `db:"status"`
	LastFourDigits   string  `db:"last_four_digits"`
	CardType         string  `db:"card_type"`
	ExpDate          string  `db:"exp_date"`
	Currency         string  `db:"currency"`
	SubTotal         float32 `gorm:"-"`
}
