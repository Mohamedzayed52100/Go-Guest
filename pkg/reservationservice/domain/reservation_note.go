package domain

import (
	"time"

	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

type ReservationNote struct {
	ID            int32            `db:"id"`
	Reservation   *Reservation     `gorm:"-"`
	ReservationID int32            `db:"reservation_id"`
	Description   string           `db:"description"`
	CreatorID     int32            `db:"creator_id"`
	Creator       *userDomain.User `gorm:"-"`
	CreatedAt     time.Time        `db:"created_at"`
	UpdatedAt     time.Time        `db:"updated_at"`
}
