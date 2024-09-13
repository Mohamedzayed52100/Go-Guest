package domain

import (
	"encoding/json"
	"time"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

type ReservationLog struct {
	ID            int32        `db:"id"`
	ReservationID int32        `db:"reservation_id"`
	CreatorID     int32        `db:"creator_id"`
	Creator       *domain.User `gorm:"-"`
	MadeBy        string       `db:"made_by"`
	FieldName     string       `db:"field_name"`
	OldValue      string       `db:"old_value"`
	NewValue      string       `db:"new_value"`
	Action        string       `db:"action"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
}

func (log *ReservationLog) ToString() string {
	result, err := json.Marshal(log)
	if err != nil {
		return ""
	}

	return string(result)
}
