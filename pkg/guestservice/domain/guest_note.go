package domain

import (
	"time"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

type GuestNote struct {
	ID          int32        `db:"id"`
	GuestID     int32        `db:"guest_id"`
	Description string       `db:"description"`
	CreatorID   int32        `db:"creator_id"`
	Creator     *domain.User `gorm:"foreignKey:CreatorID"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}
