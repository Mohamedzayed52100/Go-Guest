package domain

import (
	"time"

	tagDomain "github.com/goplaceapp/goplace-settings/pkg/guesttagservice/domain"
	"gorm.io/gorm"
)

type Guest struct {
	ID                  int32                 `db:"id"`
	FirstName           string                `db:"first_name"`
	LastName            string                `db:"last_name"`
	Email               *string               `db:"email"`
	PhoneNumber         string                `db:"phone_number"`
	Language            string                `db:"language"`
	Birthdate           *time.Time            `db:"birthdate"`
	LastVisit           *time.Time            `gorm:"-"`
	TotalVisits         int32                 `gorm:"-"`
	CurrentMood         string                `gorm:"-"`
	TotalSpent          float32               `gorm:"-"`
	TotalNoShow         int32                 `gorm:"-"`
	TotalCancel         int32                 `gorm:"-"`
	UpcomingReservation string                `gorm:"-"`
	Branches            []*GuestBranchVisit   `gorm:"-"`
	Tags                []*tagDomain.GuestTag `gorm:"-"`
	Notes               []*GuestNote          `gorm:"-"`
	Address             string                `db:"address"`
	Gender              string                `db:"gender"`
	IsPrimary           bool                  `gorm:"-"`
	CreatedAt           time.Time             `db:"created_at"`
	UpdatedAt           time.Time             `db:"updated_at"`
	DeletedAt           gorm.DeletedAt
}
