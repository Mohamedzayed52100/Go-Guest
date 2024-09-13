package repository

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	waitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationRepository) GetAvailableTimes(ctx context.Context, req *guestProto.GetAvailableTimesRequest) (*guestProto.GetAvailableTimesResponse, error) {
	var (
		resultTimes       []map[string]interface{}
		turnoverList      = []*shiftDomain.Turnover{}
		turnoverMap       = make(map[int]int32)
		shift             shiftDomain.Shift
		isOverBooking     bool
		maxTablePartySize int32
		totalReservations = []*domain.Reservation{}
		currentTime       = r.CommonRepo.ConvertToLocalTimeByBranch(ctx, time.Now(), req.GetBranchId())
		allTables         []*tableDomain.Table
	)

	if err := r.GetTenantDBConnection(ctx).
		Model(&tableDomain.Table{}).
		Select("max(max_party_size)").
		Scan(&maxTablePartySize).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	r.GetTenantDBConnection(ctx).
		Where("shift_id = ?", req.GetShiftId()).
		Find(&turnoverList)

	for _, turnover := range turnoverList {
		turnoverMap[turnover.GuestsNumber] = turnover.TurnoverTime
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Shift{}).
		Where("id = ? AND branch_id = ?", req.GetShiftId(), req.GetBranchId()).
		First(&shift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrShiftNotFound)
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	parseDate, err := time.Parse(time.DateOnly, req.GetDate())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrInvalidDate)
	}
	if !isDateInShiftRange(parseDate, shift.StartDate, shift.EndDate) {
		return &guestProto.GetAvailableTimesResponse{
			AvailableTimes: []*guestProto.AvailableTime{},
		}, nil
	}

	date, err := time.Parse(time.DateOnly, req.GetDate())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrInvalidDate)
	}

	if err := r.GetTenantDBConnection(ctx).
		Table("reservations").
		Joins("JOIN turnover ON least(10,reservations.guests_number) = turnover.guests_number").
		Joins("JOIN reservation_statuses ON reservation_statuses.id = reservations.status_id").
		Where("reservations.branch_id = ? AND "+
			"reservations.shift_id = ? AND "+
			"reservations.seating_area_id = ? AND "+
			"UPPER(reservation_statuses.type) <> UPPER(?) AND "+
			"UPPER(reservation_statuses.type) <> UPPER(?) AND "+
			"reservations.deleted_at IS NULL AND "+
			"reservations.date = ?",
			req.GetBranchId(),
			req.GetShiftId(),
			req.GetSeatingAreaId(),
			meta.Cancelled,
			meta.NoShow,
			date,
		).
		Distinct().
		Select("reservations.*").
		Scan(&totalReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	seatingAreas := r.seatingAreaClient.Client.SeatingAreaService.Repository.GetAllShiftSeatingAreas(ctx, req.GetShiftId())

	var isFound bool
	for _, sa := range seatingAreas {
		if sa.ID == req.GetSeatingAreaId() {
			isFound = true
			break
		}
	}

	if !isFound {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrSeatingAreaNotFound)
	}

	r.GetTenantDBConnection(ctx).
		Where("branch_id=?", req.GetBranchId()).
		Find(&allTables)

	for i := shift.From; ; i = i.Add(time.Duration(shift.TimeInterval * int(time.Minute))) {
		reservations := []*domain.Reservation{}
		maxPacing := 0
		waitlist := []*waitlistDomain.ReservationWaitlist{}
		fullTime := time.Date(date.Year(), date.Month(), date.Day(), i.Hour(), i.Minute(), 0, 0, time.UTC)

		for _, res := range totalReservations {
			currentTurnover := turnoverMap[int(min(10, res.GuestsNumber))]
			convertedTime, err := time.Parse(time.TimeOnly, res.Time)
			if err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}

			reservationTimestamp := time.Date(date.Year(), date.Month(), date.Day(), convertedTime.Hour(), convertedTime.Minute(), 0, 0, time.UTC)
			endTime := reservationTimestamp.Add(time.Duration(currentTurnover) * time.Minute)

			if !fullTime.Before(reservationTimestamp) && fullTime.Before(endTime) {
				reservations = append(reservations, res)
			}
		}

		if err := r.GetTenantDBConnection(ctx).
			Table("reservation_waitlists").
			Where("shift_id = ? AND seating_area_id = ?", req.GetShiftId(), req.GetSeatingAreaId()).
			Distinct().
			Find(&waitlist).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		pacing := 0
		pacingCapacity := &shiftDomain.Pacing{}
		excludedTables := []int{}
		for _, res := range reservations {
			pacing += int(res.GuestsNumber)

			var tables []*tableDomain.Table
			if err := r.GetTenantDBConnection(ctx).
				Model(&tableDomain.Table{}).
				Joins("JOIN reservation_tables ON reservation_tables.table_id = tables.id").
				Where("reservation_tables.reservation_id = ?", res.ID).Find(&tables).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}

			for _, table := range tables {
				excludedTables = append(excludedTables, table.ID)
				maxPacing += table.MaxPartySize
			}
		}

		for _, w := range waitlist {
			createdTime := r.
				CommonRepo.
				ConvertToLocalTime(ctx, w.CreatedAt).
				Add(time.Duration(w.WaitingTime))

			createdTime = time.Date(
				fullTime.Year(), fullTime.Month(), fullTime.Day(),
				createdTime.Hour(), createdTime.Minute(), 0, 0, time.UTC)

			if createdTime.Format(time.DateOnly) != req.GetDate() {
				continue
			}

			var waitlistTurnover *shiftDomain.Turnover
			if err := r.GetTenantDBConnection(ctx).
				First(&waitlistTurnover, "shift_id = ? AND guests_number = ?", w.ShiftID, w.GuestsNumber).
				Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}

			endTime := createdTime.Add(time.Duration(waitlistTurnover.TurnoverTime) * time.Minute)
			if !fullTime.Before(createdTime) && fullTime.Before(endTime) {
				pacing += int(w.GuestsNumber)
			}
		}

		if err := r.GetTenantDBConnection(ctx).
			Model(&shiftDomain.Pacing{}).
			Where("hour = ? AND "+
				"shift_id = ? AND "+
				"(seating_area_id = ? OR seating_area_id IS NULL)",
				i.Format("15:04"),
				shift.ID,
				req.GetSeatingAreaId()).
			First(&pacingCapacity).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Error(http.StatusInternalServerError, "Pacing capacity not found")
			}

			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		if !i.Equal(shift.From) && i.Hour() < shift.From.Hour() {
			fullTime = fullTime.Add(time.Hour * 24)
		}

		currentTimeAsMinutes := currentTime.Hour()*60 + currentTime.Minute()
		fullTimeAsMinutes := fullTime.Hour()*60 + fullTime.Minute()
		endTimeAsMinutes := shift.To.Hour()*60 + shift.To.Minute()

		isFullTimeAfterMidnight := fullTime.Hour() >= 0 && fullTime.Hour() < 12
		isCurrentTimeAfterMidnight := currentTime.Hour() >= 0 && currentTime.Hour() < 12
		isReqAfterCurrent := false
		isBefore := false

		if isFullTimeAfterMidnight && !isCurrentTimeAfterMidnight &&
			fullTimeAsMinutes <= endTimeAsMinutes {
			isReqAfterCurrent = true
		} else if !isFullTimeAfterMidnight && isCurrentTimeAfterMidnight {
			isBefore = true
		}
		if fullTimeAsMinutes < currentTimeAsMinutes &&
			!isReqAfterCurrent {
			isBefore = true
		}

		if fullTime.Truncate(time.Minute).After(time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, time.UTC)) {
			isBefore = false
		}

		dayBefore := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
		if r.CommonRepo.ConvertToLocalTime(ctx, time.Now()).Hour() < 12 {
			dayBefore = dayBefore.AddDate(0, 0, -1)
		}

		if !isBefore {
			available := false

			if date.After(dayBefore) && pacing < int(pacingCapacity.Capacity) {
				available = true
			} else if (req.GetPartySize() > maxTablePartySize || int32(maxPacing)+req.GetPartySize() > pacingCapacity.Capacity || pacing == int(pacingCapacity.Capacity)) && !isOverBooking {
				available = false
			}

			isTableAvailable := false
			for _, table := range allTables {
				if !contains(excludedTables, table.ID) &&
					table.MaxPartySize >= int(req.GetPartySize()) &&
					table.MinPartySize <= int(req.GetPartySize()) {
					isTableAvailable = true
					break
				}
			}

			resultTimes = append(resultTimes, map[string]interface{}{
				"time":      fullTime,
				"pacing":    pacing,
				"capacity":  pacingCapacity.Capacity,
				"available": available && isTableAvailable,
			})
		} else if !date.After(dayBefore) {
			resultTimes = append(resultTimes, map[string]interface{}{
				"time":      fullTime,
				"pacing":    pacing,
				"capacity":  pacingCapacity.Capacity,
				"available": false,
			})
		}

		if i.Hour() == shift.To.Hour() && i.Minute() == shift.To.Minute() {
			break
		}
	}

	return &guestProto.GetAvailableTimesResponse{
		AvailableTimes: converters.BuildAvailableTimesResponse(resultTimes),
	}, nil
}

func isDateInShiftRange(date, fromDate, toDate time.Time) bool {
	return (date.After(fromDate) && date.Before(toDate)) || (date.Equal(fromDate) || date.Equal(toDate))
}

func contains(list []int, item int) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}
