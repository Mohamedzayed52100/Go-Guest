package repository

import (
	"context"
	"errors"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	extDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationWidgetRepository) CreateWidgetReservation(ctx context.Context, req *guestProto.CreateWidgetReservationRequest) (*guestProto.CreateWidgetReservationResponse, error) {
	var (
		bookedStatusId   int32
		visitors         []int32
		reservationModel *domain.Reservation
		shiftId          int32
	)

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationStatus{}).
		Where("name = ? AND branch_id = ?", meta.Booked, req.GetBranchId()).
		Select("id").
		Scan(&bookedStatusId).Error; err != nil {
		return nil, status.Error(http.StatusBadRequest, "No status found")
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Shift{}).
		Joins("LEFT JOIN shift_seating_areas_assignments on shift_seating_areas_assignments.shift_id = shifts.id").
		Where(`shifts.branch_id = ? AND 
			shift_seating_areas_assignments.seating_area_id = ? AND 
			((CAST(? AS time) between CAST(shifts.from AS time) AND CAST(shifts.to AS time)) OR 
			(CAST(shifts.from AS time) > CAST(shifts.to AS time) AND (CAST(? AS time) >= CAST(shifts.from AS time) OR CAST(? AS time) <= CAST(shifts.to AS time)))) AND 
			(?::date between shifts.start_date AND shifts.end_date) AND 
			(shifts.days_to_repeat <> '' AND ? = ANY(string_to_array(shifts.days_to_repeat, ',')))`,
			req.GetBranchId(),
			req.GetSeatingAreaId(),
			req.GetTime(),
			req.GetTime(),
			req.GetTime(),
			req.GetDate(),
			utils.ExtractWeekday(req.GetDate()),
		).
		Select("shifts.id").
		Scan(&shiftId).Error; err != nil {
		return nil, status.Error(http.StatusBadRequest, "No shifts found")
	}

	reservationTable := &domain.ReservationTable{}
	table := &tableDomain.Table{
		MaxPartySize: 1e9,
	}

	var primaryGuest *guestDomain.Guest
	for _, guest := range req.GetGuests() {
		guestModel := guestDomain.Guest{
			FirstName:   guest.FirstName,
			LastName:    guest.LastName,
			PhoneNumber: guest.PhoneNumber,
			Language:    "Arabic",
		}

		if guest.BirthDate != "" {
			birthdate, err := time.Parse(time.DateOnly, guest.BirthDate)
			if err != nil {
				return nil, err
			}

			guestModel.Birthdate = &birthdate
		}

		if guest.Email != "" {
			guestModel.Email = &guest.Email
		}

		if err := r.GetTenantDBConnection(ctx).FirstOrCreate(&guestModel,
			"phone_number = ?", guest.PhoneNumber,
		).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, "Error creating guest")
		}

		if guest.Primary {
			primaryGuest = &guestModel
		} else {
			visitors = append(visitors, guestModel.ID)
		}
	}

	if primaryGuest == nil {
		return nil, errors.New("no primary guest found")
	}

	reservationModel = &domain.Reservation{
		GuestID:       primaryGuest.ID,
		BranchID:      req.GetBranchId(),
		ShiftID:       shiftId,
		SeatingAreaID: req.GetSeatingAreaId(),
		StatusID:      bookedStatusId,
		GuestsNumber:  req.GetGuestsNumber(),
		ReservedVia:   req.GetReservedVia(),
	}

	if req.GetSpecialOccasionId() != 0 {
		reservationModel.SpecialOccasionID = &req.SpecialOccasionId
	}

	var fullTime time.Time
	if reservationTime, err := time.Parse("15:04", req.GetTime()); err == nil {
		reservationModel.Time = reservationTime.Format("15:04")
		fullTime = time.Date(0, 0, 0, reservationTime.Hour(), reservationTime.Minute(), 0, 0, time.UTC)
	} else {
		logger.Default().Error(err)
		return nil, err
	}

	if reservationDate, err := time.Parse(time.DateOnly, req.GetDate()); err == nil {
		reservationModel.Date = reservationDate
		fullTime = time.Date(reservationDate.Year(), reservationDate.Month(), reservationDate.Day(), fullTime.Hour(), fullTime.Minute(), 0, 0, time.UTC)
	} else {
		logger.Default().Error(err)
		return nil, err
	}

	shift, err := r.shiftClient.Client.ShiftService.Repository.GetAllShiftData(ctx, shiftId)
	if err != nil {
		return nil, status.Error(http.StatusNotFound, "No shift found")
	}

	availableTimes, err := r.reservationRepository.GetAvailableTimes(ctx, &guestProto.GetAvailableTimesRequest{
		BranchId:      req.GetBranchId(),
		ShiftId:       shiftId,
		Date:          fullTime.Format(time.DateOnly),
		PartySize:     req.GetGuestsNumber(),
		SeatingAreaId: req.GetSeatingAreaId(),
	})
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Error getting available times")
	}

	var turnover *shiftDomain.Turnover
	for _, t := range shift.Turnover {
		if int(t.GuestsNumber) == min(10, int(req.GetGuestsNumber())) {
			turnover = t
			break
		}
	}
	if turnover == nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrMaxPartySize)
	}

	isFound := false
	for _, i := range availableTimes.GetAvailableTimes() {
		currentTime := time.Date(fullTime.Year(), fullTime.Month(), fullTime.Day(), i.GetTime().AsTime().Hour(), i.GetTime().AsTime().Minute(), 0, 0, time.UTC)
		if currentTime.Equal(fullTime) && i.Available && i.Pacing+req.GetGuestsNumber() <= i.Capacity {
			isFound = true
			break
		}
	}

	convertedTime := r.commonRepository.ConvertToLocalTime(ctx, time.Now())
	currentDay := time.Date(fullTime.Year(), fullTime.Month(), fullTime.Day(), 0, 0, 0, 0, time.UTC)
	dayBefore := time.Date(convertedTime.Year(), convertedTime.Month(), convertedTime.Day()-1, 0, 0, 0, 0, time.UTC)

	if !isFound {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
	}

	isTimeAfterMidnight := fullTime.Hour() >= 0 && fullTime.Hour() < 12 && shift.From.After(shift.To)

	if (!isTimeAfterMidnight || convertedTime.Hour()*60+convertedTime.Minute() > fullTime.Hour()*60+fullTime.Minute()+shift.TimeInterval) && currentDay.Equal(dayBefore) {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrInvalidTime)
	}

	if shift.From.After(shift.To) &&
		fullTime.Hour() >= 0 && shift.From.Hour()*60+shift.From.Minute() > fullTime.Hour()*60+fullTime.Minute() &&
		fullTime.Hour()*60+fullTime.Minute() <= (shift.To.Hour()+2)*60+shift.To.Minute() {
		fullTime = fullTime.AddDate(0, 0, -1)
	}

	if currentDay.Before(dayBefore) {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrInvalidDate)
	}

	pacing, err := r.shiftClient.Client.ShiftService.Repository.GetPacingByTime(
		ctx,
		shift.ID,
		fullTime,
		req.GetBranchId(),
		req.GetSeatingAreaId(),
	)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var capacity int32
	if err := r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Pacing{}).
		Where("shift_id = ? AND seating_area_id = ? AND hour = ?",
			shift.ID,
			req.GetSeatingAreaId(),
			fullTime.Format("15:04"),
		).
		Order("id desc").
		Select("capacity").
		Scan(&capacity).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	tables, err := r.reservationRepository.CreateReservationTables(ctx,
		fullTime,
		fullTime,
		fullTime.Add(time.Duration(turnover.TurnoverTime)*time.Minute),
		int(req.GetGuestsNumber()),
		int(req.GetBranchId()),
	)
	if err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
	}

	var minPartySize int
	r.GetTenantDBConnection(ctx).
		Model(&tableDomain.Table{}).
		Where("seating_area_id = ?", req.GetSeatingAreaId()).
		Select("MIN(max_party_size)").
		Scan(&minPartySize)

	if len(tables) == 0 && minPartySize <= int(req.GetGuestsNumber()) {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrMaxPartySize)
	}

	sameSeatingArea := false
	for _, t := range tables {
		if t.SeatingAreaID == int(req.GetSeatingAreaId()) {
			if t.MaxPartySize >= int(req.GetGuestsNumber()) &&
				t.MinPartySize <= int(req.GetGuestsNumber()) &&
				table.MaxPartySize >= t.MaxPartySize {
				table = t
			}
			sameSeatingArea = true
		}
	}

	if !sameSeatingArea && minPartySize <= int(req.GetGuestsNumber()) {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrNoTablesAvailableInSeatingArea)
	}

	if table.MaxPartySize == 1e9 && minPartySize <= int(req.GetGuestsNumber()) {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrMaxPartySize)
	}

	if capacity < pacing+req.GetGuestsNumber() {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
	}

	if err := r.GetTenantDBConnection(ctx).Create(reservationModel).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, guestId := range visitors {
		if err := r.GetTenantDBConnection(ctx).Table("reservation_visitors").Create(map[string]interface{}{
			"reservation_id": reservationModel.ID,
			"guest_id":       guestId,
		}).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	if table.CombinedTables != nil {
		for _, t := range utils.ConvertStringToArrayBySeparator(*table.CombinedTables, ",") {
			tableId, _ := strconv.ParseInt(t, 10, 32)
			reservationTable = &domain.ReservationTable{
				ReservationID: int(reservationModel.ID),
				TableID:       int(tableId),
			}
			if err := r.GetTenantDBConnection(ctx).Create(reservationTable).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}
	} else if table.ID != 0 {
		reservationTable = &domain.ReservationTable{
			ReservationID: int(reservationModel.ID),
			TableID:       table.ID,
		}
		if err := r.GetTenantDBConnection(ctx).Create(reservationTable).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	if req.GetNote() != "" {
		reservationNote := &extDomain.ReservationNote{
			ReservationID: reservationModel.ID,
			Description:   req.GetNote(),
		}
		if err := r.GetTenantDBConnection(ctx).Create(reservationNote).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &guestProto.CreateWidgetReservationResponse{
		Code:    http.StatusOK,
		Message: "Reservation created successfully",
	}, nil
}
