package domain

import (
	"encoding/json"
	"time"
)

type ReservationTagsAssignment struct {
	ID            int32     `db:"id"`
	TagID         int32     `db:"tag_id"`
	ReservationID int32     `db:"reservation_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (tag *ReservationTagsAssignment) ToString() string {
	result, err := json.Marshal(tag)
	if err != nil {
		return ""
	}

	return string(result)
}
