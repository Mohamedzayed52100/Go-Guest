package domain

import "time"

type PaymentRequestContact struct {
	ID               string    `db:"id"`
	PaymentRequestID int32     `db:"payment_request_id"`
	ContactID        int32     `db:"contact_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
