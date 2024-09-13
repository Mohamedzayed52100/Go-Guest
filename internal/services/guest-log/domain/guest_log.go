package domain

import (
	"time"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

type GuestLog struct {
	ID        int32        `db:"id"`
	GuestID   int32        `db:"guest_id"`
	CreatorID int32        `db:"creator_id"`
	Creator   *domain.User `gorm:"-"`
	MadeBy    string       `db:"made_by"`
	FieldName string       `db:"field_name"`
	OldValue  string       `db:"old_value"`
	NewValue  string       `db:"new_value"`
	Action    string       `db:"action"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}
