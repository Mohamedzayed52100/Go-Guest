package domain

import (
	"time"

	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
)

type ReservationFeedback struct {
	ID            int32                          `db:"id"`
	Guest         *domain.Guest                  `gorm:"-"`
	ReservationID int32                          `db:"reservation_id"`
	Reservation   *reservationDomain.Reservation `gorm:"-"`
	StatusID      int32                          `db:"status_id"`
	Status        string                         `gorm:"-"`
	SolutionID    int32                          `db:"solution_id"`
	Solution      *ReservationFeedbackSolution   `gorm:"-"`
	Sections      []*ReservationFeedbackSection  `gorm:"-"`
	Rate          int32                          `db:"rate"`
	Description   string                         `db:"description"`
	CreatedAt     time.Time                      `db:"created_at"`
	UpdatedAt     time.Time                      `db:"updated_at"`
}

type SimpleReservationFeedback struct {
	ID          int32     `db:"id"`
	Rate        int32     `db:"rate"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (SimpleReservationFeedback) TableName() string {
	return "reservation_feedbacks"
}
