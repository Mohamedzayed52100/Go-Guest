package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/adapters/converters"
	reservationWaitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	reservationDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	externalWaitlistDomain "github.com/goplaceapp/goplace-guest/pkg/waitlistservice/domain"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationWaitListRepository) CreateWaitingReservation(ctx context.Context, req *guestProto.CreateWaitingReservationRequest) (*guestProto.CreateWaitingReservationResponse, error) {
	var (
		waitingReservation *reservationWaitlistDomain.ReservationWaitlist
	)

	if err := r.GetTenantDBConnection(ctx).
		First(&shiftDomain.Shift{}, "branch_id = ? AND id = ?",
			r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx),
			req.GetParams().GetShiftId(),
		).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrShiftNotFound)
	}

	getLoggedInUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusNotFound, "User not found")
	}

	permissions := []string{}
	r.GetTenantDBConnection(ctx).
		Model(&permissions).
		Joins("JOIN role_permission_assignments ON role_permission_assignments.permission_id = permissions.id").
		Where("role_permission_assignments.role_id = ?", getLoggedInUser.RoleID).
		Table("permissions").
		Select("name").
		Scan(&permissions)

	var branchId int32
	if err := r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Shift{}).
		Where("id = ?", req.GetParams().GetShiftId()).
		Select("branch_id").
		Scan(&branchId).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrShiftNotFound)
	}

	waitingReservation = &reservationWaitlistDomain.ReservationWaitlist{
		GuestID:       req.GetParams().GetGuestId(),
		SeatingAreaID: req.GetParams().GetSeatingAreaId(),
		ShiftID:       req.GetParams().GetShiftId(),
		GuestsNumber:  req.GetParams().GetGuestsNumber(),
		WaitingTime:   req.GetParams().GetWaitingTime(),
		BranchID:      branchId,
		NoteID:        nil,
		CreatorID:     getLoggedInUser.ID,
		Date:          req.GetParams().GetDate(),
	}

	if req.GetParams().GetNoteId() != 0 {
		waitingReservation.NoteID = &req.Params.NoteId
	}

	if err := r.GetTenantDBConnection(ctx).Create(&waitingReservation).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(http.StatusConflict, errorhelper.ErrDuplicateGuestInWaitlist)
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, i := range req.GetParams().GetTags() {
		tag := &reservationDomain.ReservationTag{}

		tagAssignment := &externalWaitlistDomain.ReservationWaitlistTagsAssignment{
			TagID:         i.GetId(),
			ReservationID: waitingReservation.ID,
		}

		if err := r.GetTenantDBConnection(ctx).
			First(&tag, "id = ? AND category_id = ?", i.GetId(), i.GetCategoryId()).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		if err := r.GetTenantDBConnection(ctx).
			Create(&tagAssignment).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		if err := r.GetTenantDBConnection(ctx).
			First(&tag.Category, "id = ?", tag.CategoryID).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		waitingReservation.Tags = append(waitingReservation.Tags, tag)
	}

	waitingReservation, err = r.GetAllReservationWaitlistData(ctx, waitingReservation.ID)
	if err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	if _, err := r.CreateReservationWaitListLogs(ctx, &reservationWaitlistDomain.ReservationWaitlistLog{
		ReservationWaitlistID: waitingReservation.ID,
		FieldName:             "reservation-waitlist",
		Action:                "create",
	}); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	result := &guestProto.CreateWaitingReservationResponse{
		Result: converters.BuildReservationWaitListResponse(waitingReservation),
	}

	return result, nil
}

func (r *ReservationWaitListRepository) GetWaitingReservationData(ctx context.Context, id int32) (*reservationWaitlistDomain.ReservationWaitlist, error) {
	var (
		result *reservationWaitlistDomain.ReservationWaitlist
		err    error
		note   *reservationWaitlistDomain.ReservationWaitlistNote
	)

	if err := r.GetTenantDBConnection(ctx).
		First(&result, "id = ?", id).
		Error; err != nil {
		return nil, err
	}

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
