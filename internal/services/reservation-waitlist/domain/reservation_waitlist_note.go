package domain

import (
	"time"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

type ReservationWaitlistNote struct {
	ID          int32        `db:"id"`
	Description string       `db:"description"`
	CreatorID   int32        `db:"creator_id"`
	Creator     *domain.User `gorm:"-"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}
