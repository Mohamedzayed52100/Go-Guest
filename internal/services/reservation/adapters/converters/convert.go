package converters

import (
	"log"
	"time"

	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	specialOccasionDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/domain"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	extDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	seatingAreaDomain "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildReservationResponse(req *domain2.Reservation) *guestProto.Reservation {
	parsedTime, err := time.Parse(time.TimeOnly, req.Time)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return nil
	}

	res := &guestProto.Reservation{
		Id:             req.ID,
		Guests:         converters.BuildAllReservationGuestsResponse(req.Guests),
		ReservationRef: req.ReservationRef,
		Branch: &guestProto.ReservationBranch{
			Id:   int32(req.Branch.ID),
			Name: req.Branch.Name,
		},
		Shift: &guestProto.ReservationShift{
			Id:   int32(req.Shift.ID),
			Name: req.Shift.Name,
		},
		Tables:       BuildTablesResponse(req.Tables),
		GuestsNumber: req.GuestsNumber,
		Date:         timestamppb.New(req.Date),
		Time:         timestamppb.New(parsedTime),
		ReservedVia:  req.ReservedVia,
		Status:       BuildReservationStatusResponse(req.Status),
		TotalSpent:   req.TotalSpent,
		CreatedAt:    timestamppb.New(req.CreatedAt),
		UpdatedAt:    timestamppb.New(req.UpdatedAt),
	}

	if req.SeatedGuests == 0 || req.SeatedGuests == req.GuestsNumber {
		res.SeatedGuests = 0
	} else {
		res.SeatedGuests = req.SeatedGuests
	}

	if req.CheckIn != nil {
		res.CheckIn = timestamppb.New(*req.CheckIn)
	}

	if req.CheckOut != nil {
		res.CheckOut = timestamppb.New(*req.CheckOut)
	}

	if req.Tags != nil {
		res.Tags = BuildReservationTagsProto(req.Tags)
	} else {
		res.Tags = nil
	}

	if req.SpecialOccasionID != nil {
		res.SpecialOccasion = BuildReservationSpecialOccasionProto(req.SpecialOccasion)
	} else {
		res.SpecialOccasion = nil
	}

	if req.Note != nil {
		res.Note = BuildReservationNoteProto(req.Note)
	} else {
		res.Note = nil
	}

	if req.Feedback != nil {
		res.Feedback = BuildShortReservationProto(req.Feedback)
	} else {
		res.Feedback = nil
	}

	if req.CreatorID != 0 {
		res.Creator = BuildCreatorProto(req.Creator)
	} else {
		res.Creator = nil
	}

	if req.Tables != nil {
		res.Tables = BuildTablesResponse(req.Tables)
	} else {
		res.Tables = nil
	}

	if req.Payment != nil {
		res.Payment = &guestProto.ReservationPayment{
			Status:      req.Payment.Status,
			TotalPaid:   req.Payment.TotalPaid,
			TotalUnPaid: req.Payment.TotalUnPaid,
		}
	} else {
		res.Payment = nil
	}

	return res
}

func BuildReservationShortResponse(req *domain2.Reservation) *guestProto.ReservationShort {
	parsedTime, err := time.Parse(time.TimeOnly, req.Time)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return nil
	}

	res := &guestProto.ReservationShort{
		Id:           req.ID,
		GuestsNumber: req.GuestsNumber,
		Date:         timestamppb.New(req.Date),
		Time:         timestamppb.New(parsedTime),
		ReservedVia:  req.ReservedVia,
		Status:       BuildReservationStatusResponse(req.Status),
		CreatedAt:    timestamppb.New(req.CreatedAt),
		UpdatedAt:    timestamppb.New(req.UpdatedAt),
	}

	if req.SeatedGuests == 0 || req.SeatedGuests == req.GuestsNumber {
		res.SeatedGuests = 0
	} else {
		res.SeatedGuests = req.SeatedGuests
	}

	if req.CheckIn != nil {
		res.CheckIn = timestamppb.New(*req.CheckIn)
	}

	if req.CheckOut != nil {
		res.CheckOut = timestamppb.New(*req.CheckOut)
	}

	if req.SpecialOccasionID != nil {
		res.SpecialOccasion = BuildReservationSpecialOccasionProto(req.SpecialOccasion)
	} else {
		res.SpecialOccasion = nil
	}

	return res
}

func BuildAllReservationsResponse(reservations []*domain2.Reservation) []*guestProto.Reservation {
	result := make([]*guestProto.Reservation, 0)
	for _, r := range reservations {
		result = append(result, BuildReservationResponse(r))
	}

	return result
}

func BuildSeatingAreaResponse(area *seatingAreaDomain.SeatingArea) *guestProto.SeatingArea {
	return &guestProto.SeatingArea{
		Id:   area.ID,
		Name: area.Name,
	}
}

func BuildTablesResponse(tables []*tableDomain.Table) []*guestProto.Table {
	result := []*guestProto.Table{}
	for _, t := range tables {
		result = append(result, &guestProto.Table{
			Id:           int32(t.ID),
			SeatingArea:  BuildSeatingAreaResponse(t.SeatingArea),
			TableNumber:  t.TableNumber,
			PosNumber:    int32(t.PosNumber),
			MinPartySize: int32(t.MinPartySize),
			MaxPartySize: int32(t.MaxPartySize),
			CreatedAt:    timestamppb.New(t.CreatedAt),
			UpdatedAt:    timestamppb.New(t.UpdatedAt),
		})
	}
	return result
}

func BuildReservationStatusResponse(status *domain2.ReservationStatus) *guestProto.ReservationStatus {
	return &guestProto.ReservationStatus{
		Id:        int32(status.ID),
		Name:      status.Name,
		Category:  status.Category,
		Type:      status.Type,
		Color:     status.Color,
		Icon:      status.Icon,
		CreatedAt: timestamppb.New(status.CreatedAt),
		UpdatedAt: timestamppb.New(status.UpdatedAt),
	}
}

func BuildReservationTagsProto(tags []*domain.ReservationTag) []*guestProto.Tag {
	result := make([]*guestProto.Tag, 0)
	for _, tag := range tags {
		result = append(result, &guestProto.Tag{
			Id:       int32(tag.ID),
			Name:     tag.Name,
			Category: BuildReservationTagCategoryProto(tag.Category),
		})
	}
	return result
}

func BuildReservationTagCategoryProto(category *domain.ReservationTagCategory) *guestProto.TagCategory {
	return &guestProto.TagCategory{
		Id:             category.ID,
		Name:           category.Name,
		Color:          category.Color,
		Classification: category.Classification,
		OrderIndex:     category.OrderIndex,
	}
}

func BuildReservationSpecialOccasionProto(occasion *specialOccasionDomain.SpecialOccasion) *guestProto.ReservationSpecialOccasion {
	return &guestProto.ReservationSpecialOccasion{
		Id:        int32(occasion.ID),
		Name:      occasion.Name,
		Color:     occasion.Color,
		Icon:      occasion.Icon,
		CreatedAt: timestamppb.New(occasion.CreatedAt),
		UpdatedAt: timestamppb.New(occasion.UpdatedAt),
	}
}

func BuildReservationNoteProto(note *extDomain.ReservationNote) *guestProto.ReservationNote {
	result := &guestProto.ReservationNote{
		Id:          note.ID,
		Description: note.Description,
		CreatedAt:   timestamppb.New(note.CreatedAt),
		UpdatedAt:   timestamppb.New(note.UpdatedAt),
	}

	if note.Creator != nil {
		result.Creator = BuildCreatorProto(note.Creator)
	}

	if note.Reservation != nil {
		result.Reservation = BuildReservationResponse(note.Reservation)
	}

	return result
}

func BuildCreatorProto(creator *userDomain.User) *guestProto.CreatorProfile {
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

func BuildShortReservationProto(feedback *domain.SimpleReservationFeedback) *guestProto.ReservationFeedbackShort {
	return &guestProto.ReservationFeedbackShort{
		Id:          feedback.ID,
		Rate:        feedback.Rate,
		Description: feedback.Description,
		CreatedAt:   timestamppb.New(feedback.CreatedAt),
	}
}
