package converters

import (
	"fmt"
	"log"
	"time"

	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationFeedbackDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	specialOccasionDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/domain"
	reservationDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	extDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	floorPlanDomain "github.com/goplaceapp/goplace-settings/pkg/floorplanservice/domain"
	guestTagDomain "github.com/goplaceapp/goplace-settings/pkg/guesttagservice/domain"
	seatingAreaDomain "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/domain"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Guest Branch Visits
func BuildGuestBranchVisitsResponse(branches []*domain.GuestBranchVisit) []*guestProto.GuestBranchVisits {
	var result []*guestProto.GuestBranchVisits
	for _, branch := range branches {
		result = append(result, &guestProto.GuestBranchVisits{
			BranchName: branch.Name,
			Visits:     branch.Visits,
		})
	}
	return result
}

// Guest Tags
func BuildTagsResponse(tags []*guestTagDomain.GuestTag) []*guestProto.Tag {
	var result []*guestProto.Tag
	for _, tag := range tags {
		result = append(result, &guestProto.Tag{
			Id:        tag.ID,
			Name:      tag.Name,
			Category:  BuildTagCategoryResponse(tag.Category),
			CreatedAt: timestamppb.New(tag.CreatedAt),
			UpdatedAt: timestamppb.New(tag.UpdatedAt),
		})
	}
	return result
}

// Tag Category
func BuildTagCategoryResponse(category *guestTagDomain.GuestTagCategory) *guestProto.TagCategory {
	return &guestProto.TagCategory{
		Id:             int32(category.ID),
		Name:           category.Name,
		Color:          category.Color,
		Classification: category.Classification,
		OrderIndex:     int32(category.OrderIndex),
	}
}

// Guest Notes
func BuildGuestNotesResponse(notes []*domain.GuestNote) []*guestProto.GuestNote {
	var result []*guestProto.GuestNote
	for _, note := range notes {
		result = append(result, BuildGuestNoteResponse(note))
	}
	return result
}

func BuildGuestNoteResponse(note *domain.GuestNote) *guestProto.GuestNote {
	return &guestProto.GuestNote{
		Id:          note.ID,
		GuestId:     note.GuestID,
		Description: note.Description,
		Creator:     BuildCreatorResponse(note.Creator),
		CreatedAt:   timestamppb.New(note.CreatedAt),
		UpdatedAt:   timestamppb.New(note.UpdatedAt),
	}
}

func BuildCreatorResponse(creator *userDomain.User) *guestProto.CreatorProfile {
	return &guestProto.CreatorProfile{
		Id:          creator.ID,
		FirstName:   creator.FirstName,
		LastName:    creator.LastName,
		PhoneNumber: creator.PhoneNumber,
		Avatar:      creator.Avatar,
		Email:       creator.Email,
		Role:        creator.Role.DisplayName,
	}
}

// Guest
func BuildGuestResponse(proto *domain.Guest) *guestProto.Guest {
	res := &guestProto.Guest{
		Id:                  proto.ID,
		FirstName:           proto.FirstName,
		LastName:            proto.LastName,
		PhoneNumber:         proto.PhoneNumber,
		Language:            proto.Language,
		TotalVisits:         proto.TotalVisits,
		CurrentMood:         proto.CurrentMood,
		TotalSpent:          proto.TotalSpent,
		TotalNoShow:         proto.TotalNoShow,
		TotalCancel:         proto.TotalCancel,
		UpcomingReservation: proto.UpcomingReservation,
		Branches:            BuildGuestBranchVisitsResponse(proto.Branches),
		Tags:                BuildTagsResponse(proto.Tags),
		Notes:               BuildGuestNotesResponse(proto.Notes),
		Gender:              proto.Gender,
		CreatedAt:           timestamppb.New(proto.CreatedAt),
		UpdatedAt:           timestamppb.New(proto.UpdatedAt),
	}

	if proto.Birthdate != nil {
		res.BirthDate = timestamppb.New(*proto.Birthdate)
	} else {
		res.BirthDate = nil
	}

	if proto.LastVisit != nil {
		res.LastVisit = timestamppb.New(*proto.LastVisit)
	} else {
		res.LastVisit = nil
	}

	if proto.Email != nil {
		res.Email = *proto.Email
	}

	return res
}

func BuildGuestShortResponse(proto *domain.Guest) *guestProto.GuestShort {
	return &guestProto.GuestShort{
		Id:          proto.ID,
		PhoneNumber: proto.PhoneNumber,
		FirstName:   proto.FirstName,
		LastName:    proto.LastName,
	}
}

func BuildAllGuestsResponse(proto []*domain.Guest) []*guestProto.Guest {
	var guests []*guestProto.Guest

	for _, guest := range proto {
		guests = append(guests, BuildGuestResponse(guest))
	}

	return guests
}

// Staff
func BuildStaffResponse(proto []*shiftDomain.Staff) []*guestProto.Staff {
	var result []*guestProto.Staff
	for _, staff := range proto {
		result = append(result, &guestProto.Staff{
			Id:          int32(staff.ID),
			CastId:      int32(staff.CastID),
			Name:        staff.Name,
			Role:        staff.Role,
			PhoneNumber: staff.PhoneNumber,
			CreatedAt:   timestamppb.New(staff.CreatedAt),
			UpdatedAt:   timestamppb.New(staff.UpdatedAt),
		})
	}
	return result
}

// Cast
func BuildCastResponse(cast *shiftDomain.Cast) *guestProto.Cast {
	return &guestProto.Cast{
		Id:        int32(cast.ID),
		Staff:     BuildStaffResponse(cast.Staff),
		CreatedAt: timestamppb.New(cast.CreatedAt),
		UpdatedAt: timestamppb.New(cast.UpdatedAt),
	}
}

// Branch
func BuildBranchResponse(branch *userDomain.Branch) *guestProto.Branch {
	return &guestProto.Branch{
		Id:        int32(branch.ID),
		Name:      branch.Name,
		UpdatedAt: timestamppb.New(branch.UpdatedAt),
		CreatedAt: timestamppb.New(branch.CreatedAt),
	}
}

// Seating Area
func BuildSeatingAreaResponse(area *seatingAreaDomain.SeatingArea) *guestProto.SeatingArea {
	return &guestProto.SeatingArea{
		Id:       area.ID,
		Name:     area.Name,
		BranchId: area.BranchID,
	}
}

func BuildAllSeatingAreasResponse(areas []*seatingAreaDomain.SeatingArea) []*guestProto.SeatingArea {
	result := make([]*guestProto.SeatingArea, 0)
	for _, area := range areas {
		result = append(result, BuildSeatingAreaResponse(area))
	}

	return result
}

// Tables
func BuildTablesResponse(tables []*tableDomain.Table) []*guestProto.Table {
	var result []*guestProto.Table
	for _, table := range tables {
		result = append(result, &guestProto.Table{
			Id:           int32(table.ID),
			SeatingArea:  BuildSeatingAreaResponse((*seatingAreaDomain.SeatingArea)(table.SeatingArea)),
			TableNumber:  table.TableNumber,
			PosNumber:    int32(table.PosNumber),
			MinPartySize: int32(table.MinPartySize),
			MaxPartySize: int32(table.MaxPartySize),
			CreatedAt:    timestamppb.New(table.CreatedAt),
			UpdatedAt:    timestamppb.New(table.UpdatedAt),
		})
	}
	return result
}

// Shift
func BuildShiftProto(shift *shiftDomain.Shift) *guestProto.Shift {
	return &guestProto.Shift{
		Id:           shift.ID,
		Name:         shift.Name,
		From:         timestamppb.New(shift.From),
		To:           timestamppb.New(shift.To),
		StartDate:    timestamppb.New(shift.StartDate),
		EndDate:      timestamppb.New(shift.EndDate),
		TimeInterval: int32(shift.TimeInterval),
		FloorPlan:    BuildFloorPlanProto(shift.FloorPlan),
		SeatingAreas: BuildAllSeatingAreasResponse(shift.SeatingAreas),
		CategoryId:   shift.CategoryID,
		MinGuests:    int32(shift.MinGuests),
		MaxGuests:    int32(shift.MaxGuests),
		DaysToRepeat: utils.ConvertStringToArrayBySeparator(shift.DaysToRepeat, ","),
		CreatedAt:    timestamppb.New(shift.CreatedAt),
		UpdatedAt:    timestamppb.New(shift.UpdatedAt),
	}
}

func BuildFloorPlanProto(floorPlan *floorPlanDomain.ShortFloorPlan) *guestProto.FloorPlan {
	return &guestProto.FloorPlan{
		Id:   floorPlan.ID,
		Name: floorPlan.Name,
	}
}

// Reservation Status
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

func BuildAllFeedbackSectionsResponse(sections []*reservationFeedbackDomain.ReservationFeedbackSection) []*guestProto.FeedbackSection {
	result := []*guestProto.FeedbackSection{}
	for _, section := range sections {
		result = append(result, BuildFeedbackSectionResponse(section))
	}
	return result
}

func BuildFeedbackSectionResponse(section *reservationFeedbackDomain.ReservationFeedbackSection) *guestProto.FeedbackSection {
	return &guestProto.FeedbackSection{
		Id:        section.ID,
		Name:      section.Name,
		CreatedAt: timestamppb.New(section.CreatedAt),
		UpdatedAt: timestamppb.New(section.UpdatedAt),
	}
}

// Reservation Feedback
func BuildReservationFeedbackResponse(feedback *reservationFeedbackDomain.ReservationFeedback) *guestProto.ReservationFeedback {
	res := &guestProto.ReservationFeedback{
		Id:          feedback.ID,
		Guest:       BuildGuestShortResponse(feedback.Guest),
		Reservation: BuildReservationShortResponse(feedback.Reservation),
		Status:      feedback.Status,
		Sections:    BuildAllFeedbackSectionsResponse(feedback.Sections),
		Rate:        feedback.Rate,
		Description: feedback.Description,
		CreatedAt:   timestamppb.New(feedback.CreatedAt),
		UpdatedAt:   timestamppb.New(feedback.UpdatedAt),
	}

	if feedback.Solution != nil {
		res.Solution = BuildReservationFeedbackSolutionResponse(feedback.Solution)
	} else {
		res.Solution = nil
	}

	return res
}

// Short Reservation Feedback without guest and reservation
func BuildShortReservationFeedbackResponse(feedback *reservationFeedbackDomain.SimpleReservationFeedback) *guestProto.ReservationFeedbackShort {
	return &guestProto.ReservationFeedbackShort{
		Id:          feedback.ID,
		Rate:        feedback.Rate,
		Description: feedback.Description,
		CreatedAt:   timestamppb.New(feedback.CreatedAt),
	}
}

// Reservation Feedbacks
func BuildAllReservationFeedbacksResponse(feedbacks []*reservationFeedbackDomain.ReservationFeedback) []*guestProto.ReservationFeedback {
	result := []*guestProto.ReservationFeedback{}

	for _, feedback := range feedbacks {
		result = append(result, BuildReservationFeedbackResponse(feedback))
	}

	return result
}

// Reservation Special Occasion
func BuildReservationSpecialOccasionResponse(proto *specialOccasionDomain.SpecialOccasion) *guestProto.ReservationSpecialOccasion {
	return &guestProto.ReservationSpecialOccasion{
		Id:        int32(proto.ID),
		Name:      proto.Name,
		Color:     proto.Color,
		Icon:      proto.Icon,
		CreatedAt: timestamppb.New(proto.CreatedAt),
		UpdatedAt: timestamppb.New(proto.UpdatedAt),
	}
}

func BuildReservationGuestResponse(proto *domain.Guest) *guestProto.ReservationGuest {
	res := &guestProto.ReservationGuest{
		Id:          proto.ID,
		FirstName:   proto.FirstName,
		LastName:    proto.LastName,
		PhoneNumber: proto.PhoneNumber,
		TotalVisits: proto.TotalVisits,
		TotalSpent:  proto.TotalSpent,
		TotalNoShow: proto.TotalNoShow,
		TotalCancel: proto.TotalCancel,
		IsPrimary:   proto.IsPrimary,
		Gender:      proto.Gender,
		Tags:        BuildTagsResponse(proto.Tags),
	}

	if proto.Notes != nil && len(proto.Notes) > 0 && proto.Notes[0].ID != 0 {
		res.Note = BuildGuestNoteResponse(proto.Notes[0])
	}

	return res
}

func BuildAllReservationGuestsResponse(proto []*domain.Guest) []*guestProto.ReservationGuest {
	result := make([]*guestProto.ReservationGuest, 0)
	for _, guest := range proto {
		result = append(result, BuildReservationGuestResponse(guest))
	}
	return result
}

func BuildReservationNoteResponse(note *extDomain.ReservationNote) *guestProto.ReservationNote {
	result := &guestProto.ReservationNote{
		Id:          note.ID,
		Description: note.Description,
		CreatedAt:   timestamppb.New(note.CreatedAt),
		UpdatedAt:   timestamppb.New(note.UpdatedAt),
	}

	if note.Creator != nil {
		result.Creator = BuildCreatorResponse(note.Creator)
	}
	return result
}

func BuildReservationTagCategoryResponse(category *reservationDomain.ReservationTagCategory) *guestProto.TagCategory {
	return &guestProto.TagCategory{
		Id:             category.ID,
		Name:           category.Name,
		Color:          category.Color,
		Classification: category.Classification,
		OrderIndex:     category.OrderIndex,
	}
}

func BuildReservationTagsResponse(tags []*reservationDomain.ReservationTag) []*guestProto.Tag {
	result := make([]*guestProto.Tag, 0)
	for _, tag := range tags {
		result = append(result, &guestProto.Tag{
			Id:       int32(tag.ID),
			Name:     tag.Name,
			Category: BuildReservationTagCategoryResponse(tag.Category),
		})
	}
	return result
}

// Reservation
func BuildReservationProto(req *domain2.Reservation) *guestProto.Reservation {
	parsedTime, err := time.Parse(time.TimeOnly, req.Time)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return nil
	}

	res := &guestProto.Reservation{
		Id:             req.ID,
		ReservationRef: fmt.Sprintf("%06v", req.ReservationRef),
		Guests:         BuildAllReservationGuestsResponse(req.Guests),
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
		res.Tags = BuildReservationTagsResponse(req.Tags)
	} else {
		res.Tags = nil
	}

	if req.SpecialOccasionID != nil {
		res.SpecialOccasion = BuildReservationSpecialOccasionResponse(req.SpecialOccasion)
	} else {
		res.SpecialOccasion = nil
	}

	if req.Note != nil {
		res.Note = BuildReservationNoteResponse(req.Note)
	} else {
		res.Note = nil
	}

	if req.Feedback != nil {
		res.Feedback = BuildShortReservationFeedbackResponse((*reservationFeedbackDomain.SimpleReservationFeedback)(req.Feedback))
	} else {
		res.Feedback = nil
	}

	if req.CreatorID != 0 {
		res.Creator = BuildCreatorResponse(req.Creator)
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
		Branch:       BuildBranchResponse(req.Branch),
		GuestsNumber: req.GuestsNumber,
		Date:         timestamppb.New(req.Date),
		Time:         timestamppb.New(parsedTime),
		ReservedVia:  req.ReservedVia,
		Status:       BuildReservationStatusResponse(req.Status),
		Tables:       BuildTablesResponse(req.Tables),
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
		res.SpecialOccasion = BuildReservationSpecialOccasionResponse(req.SpecialOccasion)
	} else {
		res.SpecialOccasion = nil
	}

	return res
}

func BuildAllReservationsResponse(proto []*domain2.Reservation) []*guestProto.Reservation {
	var result []*guestProto.Reservation
	for _, reservation := range proto {
		result = append(result, BuildReservationProto(reservation))
	}

	return result
}

func BuildAvailableTimesResponse(times []map[string]interface{}) []*guestProto.AvailableTime {
	result := []*guestProto.AvailableTime{}
	for _, t := range times {
		reqTime, _ := t["time"].(time.Time)
		pacing, _ := t["pacing"].(int)
		capacity, _ := t["capacity"].(int32)
		available, _ := t["available"].(bool)

		result = append(result, &guestProto.AvailableTime{
			Time:      timestamppb.New(reqTime),
			Pacing:    int32(pacing),
			Capacity:  capacity,
			Available: available,
		})
	}

	return result
}

// List of all Reservation Statuses count response
func BuildReservationStatusesCountResponse(result []map[string]interface{}) []*guestProto.ReservationStatusCount {
	var response []*guestProto.ReservationStatusCount
	for _, r := range result {
		response = append(response, &guestProto.ReservationStatusCount{
			Name:  r["name"].(string),
			Count: int32(r["count"].(int)),
		})
	}
	return response
}

// Guest statistics response
func BuildGuestStatisticsResponse(totalReservations int64, totalSpent float32, publicSatisfaction string) *guestProto.GuestStatistics {
	return &guestProto.GuestStatistics{
		TotalReservations:  int32(totalReservations),
		TotalSpent:         totalSpent,
		PublicSatisfaction: publicSatisfaction,
	}
}

func BuildGuestSpendingResponse(years map[int]map[string]float32) []*guestProto.YearSpending {
	var result []*guestProto.YearSpending
	for year, months := range years {
		var yearSpending []*guestProto.MonthSpending
		for month, totalSpent := range months {
			yearSpending = append(yearSpending, &guestProto.MonthSpending{
				Month:      month,
				TotalSpent: totalSpent,
			})
		}
		result = append(result, &guestProto.YearSpending{
			Year:   int32(year),
			Months: yearSpending,
		})
	}
	return result
}

// Guest reservation statistics response
func BuildGuestReservationStatisticsResponse(result map[string]int32) []*guestProto.GuestReservationStatistics {
	var response []*guestProto.GuestReservationStatistics
	for k, v := range result {
		response = append(response, &guestProto.GuestReservationStatistics{
			Name:  k,
			Value: v,
		})
	}
	return response
}

// All special occasions response
func BuildAllSpecialOccasionsResponse(specialOccasions []*domain2.SpecialOccasion) []*guestProto.SpecialOccasion {
	var result []*guestProto.SpecialOccasion
	for _, s := range specialOccasions {
		result = append(result, &guestProto.SpecialOccasion{
			Id:        int32(s.ID),
			Name:      s.Name,
			Color:     s.Color,
			Icon:      s.Icon,
			CreatedAt: timestamppb.New(s.CreatedAt),
			UpdatedAt: timestamppb.New(s.UpdatedAt),
		})
	}

	return result
}

func BuildReservationFeedbackSolutionResponse(solution *reservationFeedbackDomain.ReservationFeedbackSolution) *guestProto.ReservationFeedbackSolution {
	return &guestProto.ReservationFeedbackSolution{
		Id:        int32(solution.ID),
		Creator:   BuildCreatorResponse(solution.Creator),
		Solution:  solution.Solution,
		CreatedAt: timestamppb.New(solution.CreatedAt),
		UpdatedAt: timestamppb.New(solution.UpdatedAt),
	}
}

func BuildOrderItemResponse(item *reservationDomain.ReservationOrderItem) *guestProto.ReservationOrderItem {
	return &guestProto.ReservationOrderItem{
		Id:        int32(item.ID),
		ItemName:  item.ItemName,
		Cost:      float32(item.Cost),
		Quantity:  int32(item.Quantity),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
	}
}

func BuildAllOrderItemsResponse(items []*reservationDomain.ReservationOrderItem) []*guestProto.ReservationOrderItem {
	response := []*guestProto.ReservationOrderItem{}
	for _, item := range items {
		response = append(response, BuildOrderItemResponse(item))
	}
	return response
}

func BuildReservationOrderResponse(order *reservationDomain.ReservationOrder) *guestProto.ReservationOrder {
	res := &guestProto.ReservationOrder{
		Id:             int32(order.ID),
		DiscountAmount: float32(order.DiscountAmount),
		DiscountReason: order.DiscountReason,
		PrevailingTax:  float32(order.PrevailingTax),
		Tax:            float32(order.Tax),
		SubTotal:       float32(order.Subtotal),
		FinalTotal:     float32(order.FinalTotal),
		CreatedAt:      timestamppb.New(order.CreatedAt),
		UpdatedAt:      timestamppb.New(order.UpdatedAt),
	}

	if order.Items != nil {
		res.Items = BuildAllOrderItemsResponse(order.Items)
	}

	return res
}

func BuildCoverFlowResponse(coverFlow []*reservationDomain.CoverFlow) []*guestProto.CoverFlow {
	var response []*guestProto.CoverFlow
	for _, c := range coverFlow {
		response = append(response, &guestProto.CoverFlow{
			Time:         c.Time,
			Reservations: BuildCoverFlowReservationsResponse(c.Reservations),
		})
	}

	return response
}

func BuildCoverFlowReservationsResponse(reservations []*reservationDomain.CoverFlowReservation) []*guestProto.CoverFlowReservation {
	var response []*guestProto.CoverFlowReservation
	for _, r := range reservations {
		response = append(response, &guestProto.CoverFlowReservation{
			Id:           r.ID,
			GuestsNumber: r.GuestsNumber,
			Status: &guestProto.CoverFlowReservationStatus{
				Id:    r.Status.ID,
				Name:  r.Status.Name,
				Color: r.Status.Color,
				Icon:  r.Status.Icon,
			},
		})
	}
	return response
}
