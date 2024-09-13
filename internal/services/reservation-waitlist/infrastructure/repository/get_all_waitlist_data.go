package repository

import (
	"context"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
)

func (r *ReservationWaitListRepository) GetAllReservationWaitlistData(ctx context.Context, id int32) (*domain.ReservationWaitlist, error) {
	var result *domain.ReservationWaitlist

	if err := r.GetTenantDBConnection(ctx).First(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}

	var err error
	result.Guest, err = r.CommonRepo.GetAllGuestData(ctx, &guestDomain.Guest{ID: result.GuestID})
	if err != nil {
		return nil, err
	}

	result.SeatingArea, err = r.seatingAreaClient.Client.SeatingAreaService.Repository.GetSeatingAreaByID(ctx, result.SeatingAreaID)
	if err != nil {
		return nil, err
	}

	result.Shift, err = r.shiftClient.Client.ShiftService.Repository.GetAllShiftData(ctx, result.ShiftID)
	if err != nil {
		return nil, err
	}

	var note *domain.ReservationWaitlistNote
	r.GetTenantDBConnection(ctx).Order("created_at desc").First(&note, "id = ?", result.NoteID)

	if note.ID != 0 {
		loggedInUser, err := r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, note.CreatorID)
		if err != nil {
			return nil, err
		}
		note.Creator = loggedInUser
		result.Note = note
	} else {
		result.Note = nil
	}

	result.Tags, err = r.GetWaitingReservationTags(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
