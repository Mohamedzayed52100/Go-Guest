package domain

import (
	"encoding/json"
	"time"
)

type ReservationTable struct {
	ID            int       `db:"id"`
	ReservationID int       `db:"reservation_id"`
	TableID       int       `db:"table_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (ReservationTable) TableName() string {
	return "reservation_tables"
}

func (table *ReservationTable) ToString() string {
	result, err := json.Marshal(table)
	if err != nil {
		return ""
	}

	return string(result)
}
