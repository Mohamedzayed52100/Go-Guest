package repository

import (
	"errors"
	"fmt"
	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	waitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func (r *ReservationWidgetRepository) GetWidgetAvailableTimes(ctx context.Context, req *guestProto.GetWidgetAvailableTimesRequest) (*guestProto.GetWidgetAvailableTimesResponse, error) {
	var (
		resultTimes       []map[string]interface{}
		turnoverList      = []*shiftDomain.Turnover{}
		turnoverMap       = make(map[int]int32)
		shifts            []*shiftDomain.Shift
		isOverBooking     bool
		maxTablePartySize int32
		totalReservations = []*domain.Reservation{}
		currentTime       = r.commonRepository.ConvertToLocalTimeByBranch(ctx, time.Now(), req.GetBranchId())
	)

	if err := r.GetTenantDBConnection(ctx).
		Model(&tableDomain.Table{}).
		Select("max(max_party_size)").
		Scan(&maxTablePartySize).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	extractedWeekdays := extractWeekdays(req.GetFromDate(), req.GetToDate())

	r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Shift{}).
		Joins("LEFT JOIN shift_seating_areas_assignments on shift_seating_areas_assignments.shift_id = shifts.id").
		Where(`shifts.branch_id = ? AND
			shift_seating_areas_assignments.seating_area_id = ? AND
			(? BETWEEN shifts.start_date AND shifts.end_date) AND
			(? BETWEEN shifts.start_date AND shifts.end_date) AND
			shifts.days_to_repeat <> ''`,
			req.GetBranchId(),
			req.GetSeatingAreaId(),
			req.GetFromDate(),
			req.GetToDate(),
		).
		Find(&shifts)

	if len(shifts) == 0 {
		return &guestProto.GetWidgetAvailableTimesResponse{
			AvailableTimes: []*guestProto.AvailableTime{},
		}, nil
	}

	parseFromDate, err := time.Parse(time.DateOnly, req.GetFromDate())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrInvalidDate)
	}
	parseToDate, err := time.Parse(time.DateOnly, req.GetToDate())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrInvalidDate)
	}

	for _, shift := range shifts {
		if !shiftContainsWeekdays(shift.DaysToRepeat, extractedWeekdays) {
			shifts = append(shifts[:0], shifts[1:]...)
		}

		if !isDateInShiftRange(parseFromDate, parseToDate, shift.StartDate, shift.EndDate) {
			return &guestProto.GetWidgetAvailableTimesResponse{
				AvailableTimes: []*guestProto.AvailableTime{},
			}, nil
		}

		r.GetTenantDBConnection(ctx).
			Where("shift_id = ?", shift.ID).
			Find(&turnoverList)

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
				"(reservations.date = ? OR reservations.date = ?)",
				req.GetBranchId(),
				shift.ID,
				req.GetSeatingAreaId(),
				meta.Cancelled,
				meta.NoShow,
				parseFromDate.Format(time.DateOnly),
				parseToDate.Format(time.DateOnly),
			).
			Distinct().
			Select("reservations.*").
			Scan(&totalReservations).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		seatingAreas := r.seatingAreaClient.Client.SeatingAreaService.Repository.GetAllShiftSeatingAreas(ctx, shift.ID)

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
	}

	for _, turnover := range turnoverList {
		turnoverMap[turnover.GuestsNumber] = turnover.TurnoverTime
	}

	// Pre-fetch data before loops
	totalReservationsMap := map[time.Time][]*domain.Reservation{}
	tableMap := map[int][]*tableDomain.Table{}

	for _, res := range totalReservations {
		convertedTime, err := time.Parse(time.TimeOnly, res.Time)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		reservationTimestamp := time.Date(parseFromDate.Year(), parseFromDate.Month(), parseFromDate.Day(), convertedTime.Hour(), convertedTime.Minute(), 0, 0, time.UTC)
		totalReservationsMap[reservationTimestamp] = append(totalReservationsMap[reservationTimestamp], res)
		turnoverMap[int(res.GuestsNumber)] = min(10, res.GuestsNumber)
	}

	tables := []*tableDomain.Table{}
	if err := r.GetTenantDBConnection(ctx).
		Model(&tableDomain.Table{}).
		Joins("JOIN reservation_tables ON reservation_tables.table_id = tables.id").
		Find(&tables).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, table := range tables {
		if int32(table.MaxPartySize) > maxTablePartySize {
			maxTablePartySize = int32(table.MaxPartySize)
		}

		if _, ok := tableMap[int(table.ID)]; !ok {
			tableMap[int(table.ID)] = []*tableDomain.Table{}
		}

		tableMap[int(table.ID)] = append(tableMap[int(table.ID)], table)
	}

	waitlist := []*waitlistDomain.ReservationWaitlist{}
	for _, shift := range shifts {
		if err := r.GetTenantDBConnection(ctx).
			Table("reservation_waitlists").
			Where("shift_id = ? AND seating_area_id = ?", shift.ID, req.GetSeatingAreaId()).
			Distinct().
			Find(&waitlist).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	for _, shift := range shifts {
		date, err := time.Parse(time.DateOnly, req.GetFromDate())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrInvalidDate)
		}

		for i := shift.From; ; i = i.Add(time.Duration(shift.TimeInterval * int(time.Minute))) {
			reservations := []*domain.Reservation{}
			maxPacing := 0
			fullTime := time.Date(date.Year(), date.Month(), date.Day(), i.Hour(), i.Minute(), 0, 0, time.UTC)

			if resList, ok := totalReservationsMap[fullTime]; ok {
				for _, res := range resList {
					currentTurnover := turnoverMap[int(res.GuestsNumber)]
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
			}

			pacing := 0
			pacingCapacity := &shiftDomain.Pacing{}
			for _, res := range reservations {
				pacing += int(res.GuestsNumber)

				if tables, ok := tableMap[int(res.ID)]; ok {
					for _, table := range tables {
						maxPacing += table.MaxPartySize
					}
				}
			}

			for _, w := range waitlist {
				createdTime := r.commonRepository.ConvertToLocalTime(ctx, w.CreatedAt).Add(time.Duration(w.WaitingTime))
				createdTime = time.Date(fullTime.Year(), fullTime.Month(), fullTime.Day(), createdTime.Hour(), createdTime.Minute(), 0, 0, time.UTC)

				if createdTime.Format(time.DateOnly) != req.GetFromDate() {
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
				Where("hour = ? AND shift_id = ? AND (seating_area_id = ? OR seating_area_id IS NULL)",
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

			currentTimeAsMinutes := currentTime.Hour()*60 + currentTime.Minute()
			fullTimeAsMinutes := fullTime.Hour()*60 + fullTime.Minute()
			endTimeAsMinutes := shift.To.Hour()*60 + shift.To.Minute()

			isFullTimeAfterMidnight := fullTime.Hour() >= 0 && fullTime.Hour() < 12
			isCurrentTimeAfterMidnight := currentTime.Hour() >= 0 && currentTime.Hour() < 12
			isReqAfterCurrent := false
			isBefore := false

			if isFullTimeAfterMidnight && !isCurrentTimeAfterMidnight && fullTimeAsMinutes <= endTimeAsMinutes {
				isReqAfterCurrent = true
			} else if !isFullTimeAfterMidnight && isCurrentTimeAfterMidnight {
				isBefore = true
			}
			if fullTimeAsMinutes < currentTimeAsMinutes && !isReqAfterCurrent {
				isBefore = true
			}

			if fullTime.Truncate(time.Minute).After(time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, time.UTC)) {
				isBefore = false
			}

			dayBefore := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
			if r.commonRepository.ConvertToLocalTime(ctx, time.Now()).Hour() < 12 {
				dayBefore = dayBefore.AddDate(0, 0, -1)
			}

			if !isBefore {
				available := false

				if date.After(dayBefore) && pacing < int(pacingCapacity.Capacity) {
					available = true
				} else if (req.GetPartySize() > maxTablePartySize ||
					int32(maxPacing)+req.GetPartySize() > pacingCapacity.Capacity ||
					pacing == int(pacingCapacity.Capacity)) &&
					!isOverBooking {
					available = false
				}

				resultTimes = append(resultTimes, map[string]interface{}{
					"time":      fullTime,
					"pacing":    pacing,
					"capacity":  pacingCapacity.Capacity,
					"available": available,
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
	}

	return &guestProto.GetWidgetAvailableTimesResponse{
		AvailableTimes: converters.BuildAvailableTimesResponse(resultTimes),
	}, nil
}

func isDateInShiftRange(fromDate, endDate, shiftStartDate, shiftEndDate time.Time) bool {
	return (fromDate.After(shiftStartDate) && endDate.Before(shiftEndDate)) ||
		(fromDate.Equal(shiftStartDate) ||
			endDate.Equal(shiftEndDate))
}

func extractWeekdays(fromDate, toDate string) []string {
	var weekdays []string

	from, err := time.Parse(time.DateOnly, fromDate)
	if err != nil {
		fmt.Println(err)
		return weekdays
	}

	to, err := time.Parse(time.DateOnly, toDate)
	if err != nil {
		fmt.Println(err)
		return weekdays
	}

	for i := from; i.Before(to) || i.Equal(to); i = i.AddDate(0, 0, 1) {
		if contains(weekdays, i.Weekday().String()) {
			continue
		}

		weekdays = append(weekdays, i.Weekday().String())
	}

	return weekdays
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func shiftContainsWeekdays(shiftDays string, weekdays []string) bool {
	for _, day := range weekdays {
		if strings.Contains(shiftDays, day) {
			return true
		}
	}
	return false
}
