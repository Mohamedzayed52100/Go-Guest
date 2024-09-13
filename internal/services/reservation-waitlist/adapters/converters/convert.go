package converters

import (
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	reservationConvert "github.com/goplaceapp/goplace-guest/internal/services/reservation/adapters/converters"
	guestConvert "github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	seatingAreaDomain "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildWaitingReservationNoteResponse(note *domain.ReservationWaitlistNote) *guestProto.ReservationWaitlistNote {
	res := &guestProto.ReservationWaitlistNote{
		Id:          note.ID,
		Description: note.Description,
		CreatedAt:   timestamppb.New(note.CreatedAt),
		UpdatedAt:   timestamppb.New(note.UpdatedAt),
	}

	if note.Creator != nil {
		res.Creator = BuildCreatorResponse(note.Creator)
	} else {
		res.Creator = nil
	}

	return res
}

func BuildCreatorResponse(creator *userDomain.User) *guestProto.CreatorProfile {
	return &guestProto.CreatorProfile{
		Id:          creator.ID,
		FirstName:   creator.FirstName,
		LastName:    creator.LastName,
		Email:       creator.Email,
		PhoneNumber: creator.PhoneNumber,
		Avatar:      creator.Avatar,
		Role:        creator.Role.DisplayName,
	}
}

func BuildSeatingAreaResponse(area *seatingAreaDomain.SeatingArea) *guestProto.SeatingArea {
	return &guestProto.SeatingArea{
		Id:   area.ID,
		Name: area.Name,
	}
}

func BuildReservationWaitListResponse(res *domain.ReservationWaitlist) *guestProto.ReservationWaitlist {
	result := &guestProto.ReservationWaitlist{
		Id:           res.ID,
		GuestsNumber: res.GuestsNumber,
		WaitingTime:  res.WaitingTime,
		Date:         res.Date,
		BranchId:     res.BranchID,
		CreatedAt:    timestamppb.New(res.CreatedAt),
		UpdatedAt:    timestamppb.New(res.UpdatedAt),
	}

	if res.Guest != nil {
		result.Guest = guestConvert.BuildGuestResponse(res.Guest)
	}

	if res.SeatingArea != nil {
		result.SeatingArea = BuildSeatingAreaResponse(res.SeatingArea)
	}

	if res.Note != nil {
		result.Note = BuildWaitingReservationNoteResponse(res.Note)
	}

	if res.Tags != nil {
		result.Tags = reservationConvert.BuildReservationTagsProto(res.Tags)
	} else {
		result.Tags = nil
	}

	if res.Shift != nil {
		result.Shift = guestConvert.BuildShiftProto(res.Shift)
	}

	return result
}

// Build wait list response
func BuildAllReservationWaitListsResponse(res []*domain.ReservationWaitlist) []*guestProto.ReservationWaitlist {
	result := make([]*guestProto.ReservationWaitlist, 0)
	for _, r := range res {
		result = append(result, BuildReservationWaitListResponse(r))
	}

	return result
}
