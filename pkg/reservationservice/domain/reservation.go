package domain

import (
	"encoding/json"
	"time"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"

	specialOccasionDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	seatingAreaDomain "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/domain"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"gorm.io/gorm"
)

type Reservation struct {
	ID                int32                                  `db:"id"`
	ReservationRef    string                                 `db:"reservation_ref"`
	GuestID           int32                                  `db:"guest_id"`
	Guests            []*guestDomain.Guest                   `gorm:"-"`
	BranchID          int32                                  `db:"branch_id"`
	Branch            *userDomain.Branch                     `gorm:"-"`
	ShiftID           int32                                  `db:"shift_id"`
	Shift             *shiftDomain.Shift                     `gorm:"-"`
	SeatingAreaID     int32                                  `db:"seating_area_id"`
	SeatingArea       *seatingAreaDomain.SeatingArea         `gorm:"foreignkey:SeatingAreaID"`
	StatusID          int32                                  `db:"status_id"`
	Status            *ReservationStatus                     `gorm:"foreignkey:StatusID"`
	Note              *ReservationNote                       `gorm:"-"`
	Tables            []*tableDomain.Table                   `gorm:"many2many:reservation_tables;"`
	GuestsNumber      int32                                  `db:"guests_number"`
	SeatedGuests      int32                                  `db:"seated_guests"`
	Date              time.Time                              `db:"date"`
	Time              string                                 `db:"time"`
	CreationDuration  float32                                `db:"creation_duration"`
	ReservedVia       string                                 `db:"reserved_via"`
	CheckIn           *time.Time                             `db:"check_in"`
	SpecialOccasionID *int32                                 `db:"special_occasion_id"`
	SpecialOccasion   *specialOccasionDomain.SpecialOccasion `gorm:"foreignkey:SpecialOccasionID"`
	Tags              []*domain.ReservationTag               `gorm:"-"`
	Feedback          *domain.SimpleReservationFeedback      `gorm:"-"`
	CreatorID         int32                                  `db:"creator_id"`
	Creator           *userDomain.User                       `gorm:"-"`
	CheckOut          *time.Time                             `db:"check_out"`
	TotalSpent        float32                                `gorm:"-"`
	Payment           *ReservationPayment                    `gorm:"-"`
	CreatedAt         time.Time                              `db:"created_at"`
	UpdatedAt         time.Time                              `db:"updated_at"`
	DeletedAt         gorm.DeletedAt
}

func (r *Reservation) ToString() string {
	result, err := json.Marshal(r)
	if err != nil {
		return ""
	}

	return string(result)
}
