package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationLogDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	waitingReservationDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	externalWaitlistDomain "github.com/goplaceapp/goplace-guest/pkg/waitlistservice/domain"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	extDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationWaitListRepository) SeatWaitingReservation(ctx context.Context, req *guestProto.SeatWaitingReservationRequest) (*guestProto.SeatWaitingReservationResponse, error) {
	var (
		logs               = []*waitingReservationDomain.ReservationWaitlistLog{}
		waitingReservation = &waitingReservationDomain.ReservationWaitlist{}
		tags               = []*externalWaitlistDomain.ReservationWaitlistTagsAssignment{}
		selectedStatus     *domain2.ReservationStatus
		oldReservation     *domain2.Reservation
	)

	if err := r.GetTenantDBConnection(ctx).
		First(&waitingReservation, "id = ? AND branch_id = ?", req.GetId(),
			r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx)).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Waiting reservation not found")
	}

	r.GetTenantDBConnection(ctx).Find(&tags, "reservation_id = ?", req.GetId())

	r.GetTenantDBConnection(ctx).Find(&logs, "reservation_waitlist_id = ?", req.GetId())

	if err := r.GetTenantDBConnection(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.GetTenantDBConnection(ctx).
			Delete(&waitingReservationDomain.ReservationWaitlistLog{}, "reservation_waitlist_id = ?", req.GetId()).
			Error; err != nil {
			logger.Default().Errorf("Error deleting reservation waitlist log: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		if err := r.GetTenantDBConnection(ctx).
			Delete(&externalWaitlistDomain.ReservationWaitlistTagsAssignment{}, "reservation_id = ?", req.GetId()).
			Error; err != nil {
			logger.Default().Errorf("Error deleting reservation waitlist tags assignment: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		if err := r.GetTenantDBConnection(ctx).
			Delete(&waitingReservationDomain.ReservationWaitlist{}, "id = ?", req.GetId()).
			Error; err != nil {
			logger.Default().Errorf("Error deleting reservation waitlist: %v", err)
			return status.Error(http.StatusInternalServerError, errorhelper.ErrReservationNotFound)
		}

		var oldNote *waitingReservationDomain.ReservationWaitlistNote
		if waitingReservation.NoteID != nil {
			if err := r.GetTenantDBConnection(ctx).
				First(&oldNote, "id =?", waitingReservation.NoteID).
				Error; err != nil {
				logger.Default().Errorf("Error getting reservation waitlist note: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}

			if err := r.GetTenantDBConnection(ctx).
				Delete(&waitingReservationDomain.ReservationWaitlistNote{}, "id = ?", waitingReservation.NoteID).
				Error; err != nil {
				logger.Default().Errorf("Error deleting reservation waitlist note: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		currentTime := r.CommonRepo.ConvertToLocalTime(ctx, time.Now())
		currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, time.UTC)

		var fullShift *shiftDomain.Shift
		if err := r.GetTenantDBConnection(ctx).
			First(&fullShift, "id = ?", waitingReservation.ShiftID).
			Error; err != nil {
			logger.Default().Errorf("Error getting shift: %v", err)
		}

		timeInterval := fullShift.TimeInterval
		reservationMinutes := currentTime.Hour()*60 + currentTime.Minute()
		remainder := int32(reservationMinutes) % int32(timeInterval)
		if remainder != 0 {
			addMinutes := int32(timeInterval) - remainder
			currentTime = currentTime.Add(time.Duration(addMinutes) * time.Minute)
		}

		if fullShift.From.After(fullShift.To) &&
			currentTime.Hour() >= 0 &&
			fullShift.From.Hour()*60+fullShift.From.Minute() > currentTime.Hour()*60+currentTime.Minute() &&
			currentTime.Hour()*60+currentTime.Minute() <= (fullShift.To.Hour()+2)*60+fullShift.To.Minute() {
			currentTime = currentTime.AddDate(0, 0, -1)
		}

		if req.GetStatus() == 0 {
			getStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.Seated, waitingReservation.BranchID)
			if err != nil {
				logger.Default().Errorf("Error getting reservation status: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}

			selectedStatus = getStatus
		} else {
			getStatus, err := r.CommonRepo.GetReservationStatusByID(ctx, req.GetStatus(), waitingReservation.BranchID)
			if err != nil {
				logger.Default().Errorf("Error getting reservation status: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}

			selectedStatus = getStatus
		}

		waitingReservationDate, err := time.Parse("2006-01-02T15:04:05Z", waitingReservation.Date)
		if err != nil {
			logger.Default().Errorf("Error parsing date: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		oldReservation = &domain2.Reservation{
			SeatingAreaID: waitingReservation.SeatingAreaID,
			StatusID:      int32(selectedStatus.ID),
			GuestID:       waitingReservation.GuestID,
			ShiftID:       waitingReservation.ShiftID,
			Tables:        []*tableDomain.Table{},
			GuestsNumber:  waitingReservation.GuestsNumber,
			BranchID:      fullShift.BranchID,
			Date:          waitingReservationDate,
			Time:          currentTime.Format(time.TimeOnly),
			ReservedVia:   "Walked in",
			CreatorID:     waitingReservation.CreatorID,
		}

		if selectedStatus.Name == meta.Seated {
			oldReservation.CheckIn = &currentTime
		}

		if req.GetType() != "" {
			oldReservation.ReservedVia = req.GetType()
		}

		primaryGuest, err := r.CommonRepo.GetAllGuestData(ctx, &guestDomain.Guest{ID: waitingReservation.GuestID})
		if err != nil {
			logger.Default().Errorf("Error getting guest: %v", err)
			return status.Error(http.StatusInternalServerError, errorhelper.ErrGuestNotFound)
		}
		primaryGuest.IsPrimary = true
		oldReservation.Guests = append(oldReservation.Guests, primaryGuest)

		oldReservation.Shift, err = r.shiftClient.Client.ShiftService.Repository.GetAllShiftData(ctx, waitingReservation.ShiftID)
		if err != nil {
			logger.Default().Errorf("Error getting shift: %v", err)
			return status.Error(http.StatusInternalServerError, errorhelper.ErrShiftNotFound)
		}

		oldReservation.Branch, err = r.userClient.Client.UserService.Repository.GetBranchByID(ctx, oldReservation.BranchID)
		if err != nil {
			logger.Default().Errorf("Error getting branch: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		oldReservation.Creator, err = r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, oldReservation.CreatorID)
		if err != nil {
			logger.Default().Errorf("Error getting user profile: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		oldReservation.Status = selectedStatus

		oldReservation.SeatingArea, err = r.seatingAreaClient.Client.SeatingAreaService.Repository.GetSeatingAreaByID(ctx, oldReservation.SeatingAreaID)
		if err != nil {
			logger.Default().Errorf("Error getting seating area: %v", err)
			return status.Error(http.StatusInternalServerError, errorhelper.ErrSeatingAreaNotFound)
		}

		var turnover *shiftDomain.Turnover
		if err := r.GetTenantDBConnection(ctx).
			First(&turnover, "shift_id = ? AND guests_number = ?",
				waitingReservation.ShiftID,
				waitingReservation.GuestsNumber).
			Error; err != nil {
			logger.Default().Errorf("Error getting turnover: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		endTime := currentTime.Add(time.Duration(turnover.TurnoverTime) * time.Minute)
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)
		oldReservation.CheckOut = &endTime

		if err := r.GetTenantDBConnection(ctx).Create(&oldReservation).Error; err != nil {
			logger.Default().Errorf("Error creating reservation: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		for _, tag := range tags {
			tagAssignment := &domain.ReservationTagsAssignment{
				TagID:         tag.TagID,
				ReservationID: oldReservation.ID,
			}

			if err := r.GetTenantDBConnection(ctx).Create(&tagAssignment).Error; err != nil {
				logger.Default().Errorf("Error creating reservation tags assignment: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		resTags, err := r.CommonRepo.GetReservationTags(ctx, oldReservation.ID)
		if err != nil {
			logger.Default().Errorf("Error getting reservation tags: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}
		oldReservation.Tags = resTags

		if selectedStatus.Name == meta.Seated {
			for _, t := range req.GetTables() {
				table := &domain2.ReservationTable{
					ReservationID: int(oldReservation.ID),
					TableID:       int(t),
				}

				if err := r.GetTenantDBConnection(ctx).Create(&table).Error; err != nil {
					logger.Default().Errorf("Error creating reservation table: %v", err)
					return status.Error(http.StatusInternalServerError, err.Error())
				}
			}

			tables, err := r.CommonRepo.GetTablesForReservation(ctx, oldReservation.ID)
			if err != nil {
				logger.Default().Errorf("Error getting tables for reservation: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}
			oldReservation.Tables = tables
		}

		if oldNote != nil {
			reservationNote := &extDomain.ReservationNote{
				ReservationID: oldReservation.ID,
				Description:   oldNote.Description,
				CreatorID:     oldNote.CreatorID,
				CreatedAt:     oldNote.CreatedAt,
				UpdatedAt:     oldNote.UpdatedAt,
			}

			loggedInUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
			if err != nil {
				logger.Default().Errorf("Error getting logged in user: %v", err)
				return status.Errorf(http.StatusNotFound, err.Error())
			}
			reservationNote.Creator = loggedInUser

			if err := r.GetTenantDBConnection(ctx).Create(&reservationNote).Error; err != nil {
				logger.Default().Errorf("Error creating reservation note: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}

			oldReservation.Note = reservationNote
		}

		for _, log := range logs {
			if _, err := r.reservationRepository.CreateReservationLogs(ctx, &reservationLogDomain.ReservationLog{
				ReservationID: oldReservation.ID,
				CreatorID:     log.CreatorID,
				MadeBy:        log.MadeBy,
				FieldName:     log.FieldName,
				OldValue:      log.OldValue,
				NewValue:      log.NewValue,
				Action:        log.Action,
				CreatedAt:     log.CreatedAt,
				UpdatedAt:     log.UpdatedAt,
			}); err != nil {
				logger.Default().Errorf("Error creating reservation log: %v", err)
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		return nil
	}); err != nil {
		logger.Default().Errorf("Error transaction: %v", err)
		return nil, err
	}

	return &guestProto.SeatWaitingReservationResponse{
		Result: converters.BuildReservationResponse(oldReservation),
	}, nil
}
