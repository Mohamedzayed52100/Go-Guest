package domain

import (
	"time"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"

	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	seatingAreaDomain "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/domain"
)

type ReservationWaitlist struct {
	ID            int32                          `db:"id"`
	GuestID       int32                          `db:"guest_id"`
	Guest         *guestDomain.Guest             `gorm:"-"`
	SeatingAreaID int32                          `db:"seating_area_id"`
	SeatingArea   *seatingAreaDomain.SeatingArea `gorm:"-"`
	ShiftID       int32                          `db:"shift_id"`
	Shift         *shiftDomain.Shift             `gorm:"-"`
	GuestsNumber  int32                          `db:"guests_number"`
	WaitingTime   int32                          `db:"waiting_time"`
	NoteID        *int32                         `db:"note_id"`
	BranchID      int32                          `db:"branch_id"`
	Note          *ReservationWaitlistNote       `gorm:"-"`
	Tags          []*domain.ReservationTag       `gorm:"-"`
	CreatorID     int32                          `db:"creator_id"`
	Creator       *userDomain.User               `gorm:"-"`
	Date          string                         `db:"date"`
	CreatedAt     time.Time                      `db:"created_at"`
	UpdatedAt     time.Time                      `db:"updated_at"`
}
