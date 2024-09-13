package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"gorm.io/gorm"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/rbac"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationLogDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) CreateReservation(ctx context.Context, req *guestProto.CreateReservationRequest) (*guestProto.CreateReservationResponse, error) {
	var (
		shift         *shiftDomain.Shift
		shiftTurnover []*shiftDomain.Turnover
	)

	userRepo := r.userClient.Client.UserService.Repository

	if err := r.GetTenantDBConnection(ctx).
		Where("id = ?", req.GetParams().GetShiftId()).
		First(&shiftDomain.Shift{}).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Shift not found")
	}

	if err := r.GetTenantDBConnection(ctx).
		Table("shifts").
		Where(`id = ? AND
		NOT EXISTS (
			SELECT 1
			FROM unnest(string_to_array(exceptions, ',')) AS exception
			WHERE exception = ?
		) AND days_to_repeat <> ''`, req.GetParams().GetShiftId(), req.GetParams().GetDate()).
		Select(`"id", "name", "from", "to", "time_interval"`).
		First(&shift).Error; err != nil {
		return nil, status.Error(http.StatusPreconditionFailed, "The shift has been edited, please reselect the new shift")
	}

	isOverBooking := userRepo.CheckAdminPinCode(ctx, req.GetParams().GetPinCode())

	currentUser, err := userRepo.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if currentUser == nil {
		isOverBooking = false
	}

	if !userRepo.CheckForBranchAccess(ctx, req.GetParams().GetBranchId()) {
		return nil, status.Error(http.StatusNotFound, "You don't have access to this branch")
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Where("guest_id = ? AND date = ? AND branch_id = ?",
			req.GetParams().GetGuestId(),
			req.GetParams().GetDate(),
			req.GetParams().GetBranchId(),
		).
		First(&reservationDomain.Reservation{}).
		Error; err == nil {
		return nil, status.Error(http.StatusConflict, "The guest had previously made a reservation for the same date")
	}

	permissions := userRepo.RoleRepository.GetAllStringPermissions(ctx)

	for _, p := range permissions {
		if p == rbac.OverbookReservation.Name {
			isOverBooking = true
			break
		}
	}

	convertedTime := r.CommonRepo.ConvertToLocalTime(ctx, time.Now())

	if err := r.GetTenantDBConnection(ctx).
		Table("turnover").
		Where("shift_id = ?", shift.ID).
		Find(&shiftTurnover).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Error getting shift turnover")
	}

	reservationDate, err := time.Parse(time.DateOnly, req.GetParams().GetDate())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrInvalidDate)
	}

	reservationTime, err := time.Parse("15:04", req.GetParams().GetTime())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Modify the reservation time to the nearest time interval
	timeInterval := shift.TimeInterval
	reservationMinutes := reservationTime.Hour()*60 + reservationTime.Minute()
	remainder := reservationMinutes % timeInterval
	if remainder != 0 {
		addMinutes := timeInterval - remainder
		reservationTime = reservationTime.Add(time.Duration(addMinutes) * time.Minute)
	}

	fullTime := time.Date(reservationDate.Year(), reservationDate.Month(), reservationDate.Day(), reservationTime.Hour(), reservationTime.Minute(), 0, 0, time.UTC)

	availableTimes, err := r.GetAvailableTimes(ctx, &guestProto.GetAvailableTimesRequest{
		BranchId:      req.GetParams().GetBranchId(),
		ShiftId:       req.GetParams().GetShiftId(),
		Date:          reservationDate.Format(time.DateOnly),
		PartySize:     req.GetParams().GetGuestsNumber(),
		SeatingAreaId: req.GetParams().GetSeatingAreaId(),
	})
	if err != nil {
		return nil, err
	}

	var turnover *shiftDomain.Turnover
	for _, t := range shiftTurnover {
		if int(t.GuestsNumber) == min(10, int(req.GetParams().GetGuestsNumber())) {
			turnover = t
			break
		}
	}
	if turnover == nil {
		return nil, status.Error(http.StatusInternalServerError, "Turnover not found")
	}

	startTime := fullTime
	endTime := fullTime.Add(time.Duration(turnover.TurnoverTime) * time.Minute)
	isFound := false
	for j := startTime; j.Before(endTime); j = j.Add(time.Duration(shift.TimeInterval) * time.Minute) {
		isFound = false
		for _, i := range availableTimes.GetAvailableTimes() {
			currentTime := time.Date(j.Year(), j.Month(), j.Day(), i.GetTime().AsTime().Hour(), i.GetTime().AsTime().Minute(), 0, 0, time.UTC)
			if currentTime.Equal(j) && i.Available && i.Pacing+req.GetParams().GetGuestsNumber() <= i.Capacity {
				isFound = true
				break
			}
		}
		if !isFound {
			break
		}
	}

	currentDay := time.Date(fullTime.Year(), fullTime.Month(), fullTime.Day(), 0, 0, 0, 0, time.UTC)
	dayBefore := time.Date(convertedTime.Year(), convertedTime.Month(), convertedTime.Day()-1, 0, 0, 0, 0, time.UTC)

	if !isFound && !isOverBooking {
		logger.Default().Error("No tables available in line 120")
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
	}

	// isTimeAfterMidnight := fullTime.Hour() >= 0 && fullTime.Hour() < 12 && shift.From.After(shift.To)

	// if (!isTimeAfterMidnight || convertedTime.Hour()*60+convertedTime.Minute() > fullTime.Hour()*60+fullTime.Minute()+shift.TimeInterval) && currentDay.Equal(dayBefore) {
	// 	return nil, status.Error(http.StatusNotFound, errorhelper.ErrInvalidTime)
	// }

	if shift.From.After(shift.To) &&
		fullTime.Hour() >= 0 && shift.From.Hour()*60+shift.From.Minute() > fullTime.Hour()*60+fullTime.Minute() &&
		fullTime.Hour()*60+fullTime.Minute() <= (shift.To.Hour()+2)*60+shift.To.Minute() {
		fullTime = fullTime.AddDate(0, 0, -1)
	}

	if currentDay.Before(dayBefore) {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrInvalidDate)
	}

	var reservationStatus *reservationDomain.ReservationStatus
	if req.GetParams().GetStatusId() != 0 {
		reservationStatus, err = r.GetReservationStatusByID(ctx, req.GetParams().GetStatusId(), req.GetParams().GetBranchId())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	} else {
		reservationStatus, err = r.GetReservationStatusByName(ctx, meta.Booked, req.GetParams().GetBranchId())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	pacing, err := r.shiftClient.Client.ShiftService.Repository.GetPacingByTime(
		ctx,
		shift.ID,
		fullTime,
		req.GetParams().GetBranchId(),
		req.GetParams().GetSeatingAreaId(),
	)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var capacity int32
	if err := r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Pacing{}).
		Where("shift_id = ? AND seating_area_id = ? AND hour = ?",
			shift.ID,
			req.GetParams().GetSeatingAreaId(),
			fullTime.Format("15:04"),
		).
		Order("id desc").
		Select("capacity").
		Scan(&capacity).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	tables, err := r.CreateReservationTables(
		ctx,
		fullTime,
		fullTime,
		fullTime.Add(time.Duration(turnover.TurnoverTime)*time.Minute),
		int(req.GetParams().GetGuestsNumber()),
		int(req.GetParams().GetBranchId()),
	)
	if err != nil && !isOverBooking {
		logger.Default().Error("No tables available in line 161")
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
	}

	logger.Default().Info("Tables: ", tables)

	var minPartySize int
	r.GetTenantDBConnection(ctx).
		Raw("SELECT MIN(max_party_size) FROM tables WHERE seating_area_id = ?", req.GetParams().GetSeatingAreaId()).
		Scan(&minPartySize)

	if len(tables) == 0 && !isOverBooking && minPartySize <= int(req.GetParams().GetGuestsNumber()) {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrMaxPartySize)
	}

	createdReservation := &reservationDomain.Reservation{
		ID:                req.GetParams().GetId(),
		GuestID:           req.GetParams().GetGuestId(),
		BranchID:          req.GetParams().GetBranchId(),
		ShiftID:           req.GetParams().GetShiftId(),
		SeatingAreaID:     req.GetParams().GetSeatingAreaId(),
		StatusID:          int32(reservationStatus.ID),
		Status:            reservationStatus,
		GuestsNumber:      req.GetParams().GetGuestsNumber(),
		Date:              reservationDate,
		Time:              reservationTime.Format(time.TimeOnly),
		CreationDuration:  req.GetParams().GetCreationDuration(),
		ReservedVia:       req.GetParams().GetReservedVia(),
		SpecialOccasionID: nil,
	}

	for {
		if err := r.GetTenantDBConnection(ctx).
			Model(&reservationDomain.Reservation{}).
			Where("reservation_ref = ?", r.GenerateReservationRef()).
			First(&reservationDomain.Reservation{}).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				createdReservation.ReservationRef = r.GenerateReservationRef()
				break
			}
		}
	}

	if currentUser != nil {
		createdReservation.CreatorID = currentUser.ID
	}

	if createdReservation.ReservedVia == "" {
		createdReservation.ReservedVia = "External"
	}

	specialOccasionID := req.GetParams().GetSpecialOccasionId()
	if specialOccasionID != 0 {
		createdReservation.SpecialOccasionID = &specialOccasionID
	} else {
		createdReservation.SpecialOccasionID = nil
	}

	reservationTable := &reservationDomain.ReservationTable{}

	if capacity < pacing+req.GetParams().GetGuestsNumber() && !isOverBooking {
		logger.Default().Error("No tables available in line 232")
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
	}

	if err := r.GetTenantDBConnection(ctx).
		Create(createdReservation).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, tag := range req.GetParams().GetTags() {
		var queryTag domain.ReservationTag

		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationTag{}).
			Where("id = ? AND category_id = ?", tag.Id, tag.CategoryId).
			First(&queryTag).
			Error; err != nil {
			continue
		}

		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationTagsAssignment{}).
			Create(&domain.ReservationTagsAssignment{
				TagID:         queryTag.ID,
				ReservationID: createdReservation.ID,
			}).Error; err != nil {
			continue
		}
	}

	if len(tables) != 0 || tables != nil {
		for _, table := range tables {
			if table.CombinedTables != nil {
				for _, t := range utils.ConvertStringToArrayBySeparator(*table.CombinedTables, ",") {
					tableId, _ := strconv.ParseInt(t, 10, 32)
					reservationTable = &reservationDomain.ReservationTable{
						ReservationID: int(createdReservation.ID),
						TableID:       int(tableId),
					}
					if err := r.GetTenantDBConnection(ctx).Create(reservationTable).Error; err != nil {
						return nil, status.Error(http.StatusInternalServerError, err.Error())
					}
				}
				break
			} else if table.ID != 0 {
				reservationTable = &reservationDomain.ReservationTable{
					ReservationID: int(createdReservation.ID),
					TableID:       table.ID,
				}

				if err := r.GetTenantDBConnection(ctx).Create(reservationTable).Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				break
			}
		}
	}

	createdReservation, err = r.CommonRepo.GetAllReservationData(ctx, createdReservation)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for i := range createdReservation.Tables {
		createdReservation.Tables[i].SeatingArea, err = r.seatingAreaClient.Client.
			SeatingAreaService.Repository.GetSeatingAreaByID(ctx, int32(createdReservation.Tables[i].SeatingAreaID))
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	if _, err := r.CreateReservationLogs(ctx, &reservationLogDomain.ReservationLog{
		ReservationID: createdReservation.ID,
		Action:        "create",
		FieldName:     "reservation",
	},
	); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if createdReservation.SpecialOccasionID != nil {
		if err := r.SendSpecialOccasionMessage(
			fmt.Sprintf("A special occasion reservation has been made for reservation ID %d with special occasion %s",
				createdReservation.ID, createdReservation.SpecialOccasion.Name),
			createdReservation.ID,
		); err != nil {
			logger.Default().Error(err.Error())
		}
	}

	if createdReservation.ReservedVia != "Direct in" && createdReservation.ReservedVia != "Walked in" {
		_, err = r.SendReservationWhatsappDetails(ctx, createdReservation)
		if err != nil {
			logger.Default().Error(err.Error())
		}
	}

	result := &guestProto.CreateReservationResponse{
		Result: converters.BuildReservationProto(createdReservation),
	}

	return result, nil
}
