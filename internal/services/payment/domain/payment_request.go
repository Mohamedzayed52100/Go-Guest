package domain

import (
	"time"
)

type PaymentRequest struct {
	ID            int32                    `db:"id"`
	Guest         *PaymentGuest            `gorm:"-"`
	ReservationID int32                    `db:"reservation_id"`
	BranchID      int32                    `db:"branch_id"`
	Branch        *PaymentBranch           `gorm:"-"`
	Delivery      string                   `db:"delivery"`
	Items         []*InvoiceItem           `gorm:"-"`
	Date          string                   `db:"date"`
	Invoice       *Invoice                 `gorm:"-"`
	Contacts      []*PaymentRequestContact `gorm:"-"`
	CreatedAt     time.Time                `db:"created_at"`
	UpdatedAt     time.Time                `db:"updated_at"`
}
