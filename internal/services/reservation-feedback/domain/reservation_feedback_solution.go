package domain

import (
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"time"
)

type ReservationFeedbackSolution struct {
	ID        int32            `db:"id"`
	CreatorID int32            `db:"creator_id"`
	Creator   *userDomain.User `gorm:"-"`
	Solution  string           `db:"solution"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
}
