package repository

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/rbac"
	commonUtils "github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationLogDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/utils"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) UpdateReservation(ctx context.Context, req *guestProto.UpdateReservationRequest) (*guestProto.UpdateReservationResponse, error) {
	var (
		err             error
		reservationTime time.Time
		reservationDate time.Time
		logs            []*reservationLogDomain.ReservationLog
		shift           *shiftDomain.Shift
		updates         = make(map[string]interface{})
		userRepo        = r.userClient.Client.UserService.Repository
	)

	isOverBooking := userRepo.CheckAdminPinCode(ctx, string(req.GetParams().GetPinCode()))

	permissions := r.roleClient.Client.RoleService.Repository.GetAllStringPermissions(ctx)

	for _, p := range permissions {
		if p == rbac.OverbookReservation.Name {
			isOverBooking = true
			break
		}
	}

	if currentUser, _ := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx); currentUser == nil {
		isOverBooking = false
	}

	currentTime := r.CommonRepo.ConvertToLocalTime(ctx, time.Now())

	currentReservation, err := r.CommonRepo.GetReservationByID(ctx, req.GetParams().GetId())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	currentReservationTime, err := time.Parse(time.TimeOnly, currentReservation.Time)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if !r.userClient.Client.UserService.Repository.CheckForBranchAccess(ctx, currentReservation.BranchID) {
		return nil, status.Error(http.StatusNotFound, "You don't have access to this branch")
	}

	// Get the turnover time for the new guests number
	if req.GetParams().GetGuestsNumber() != currentReservation.GuestsNumber && req.GetParams().GetGuestsNumber() != 0 {
		updates["guests_number"] = req.GetParams().GetGuestsNumber()

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "guests number",
			OldValue:      strconv.FormatInt(int64(currentReservation.GuestsNumber), 10),
			NewValue:      strconv.FormatInt(int64(req.GetParams().GetGuestsNumber()), 10),
		})
	} else {
		updates["guests_number"] = currentReservation.GuestsNumber
	}

	if updates["shift_id"] == nil || updates["shift_id"] == 0 {
		updates["shift_id"] = currentReservation.ShiftID
	}

	if updates["guests_number"] == nil || updates["guests_number"] == 0 {
		updates["guests_number"] = currentReservation.GuestsNumber
	}

	var turnover *shiftDomain.Turnover
	if err := r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Turnover{}).
		Where("least(?,10) = guests_number AND shift_id = ?",
			updates["guests_number"].(int32),
			updates["shift_id"].(int32)).
		First(&turnover).Error; err != nil && !isOverBooking {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrMaxPartySize)
	}

	if req.GetParams().GetStatusId() != 0 &&
		currentReservation.StatusID != req.Params.GetStatusId() {
		updates["status_id"] = req.GetParams().GetStatusId()

		getStatus, err := r.CommonRepo.GetReservationStatusByID(ctx, req.GetParams().GetStatusId(), currentReservation.BranchID)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, "Reservation status not found")
		}

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "status",
			OldValue:      currentReservation.Status.Name,
			NewValue:      getStatus.Name,
		})

		if getStatus.Type == meta.Seated {
			if currentReservation.CheckIn != nil {
				currentTime = *currentReservation.CheckIn
			}

			updates["check_in"] = currentTime

			endTime := currentTime.Add(time.Duration(turnover.TurnoverTime) * time.Minute)
			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(),
				endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)
			updates["check_out"] = endTime

			currentReservation.CheckIn = &currentTime
			currentReservation.CheckOut = &endTime
		} else if getStatus.Type == meta.Left {
			updates["check_out"] = currentTime
			currentReservation.CheckOut = &currentTime

			posIntegration, err := r.CommonRepo.GetIntegrationBySystemType(ctx, "POS System", currentReservation.BranchID)
			if err == nil {
				lCase := cases.Lower(language.English).String(posIntegration.SystemName)
				switch lCase {
				case "revel":
					if currentReservation.Tables == nil || len(currentReservation.Tables) == 0 {
						break
					}

					orderDetails, err := r.GetRevelOrderDetails(currentReservation.Tables[0].PosNumber, posIntegration, currentReservation, currentTime)
					if err != nil {
						break
					}

					if len(orderDetails) > 0 {
						reservationOrder, err := r.CreateOrUpdateReservationOrderFromRevel(ctx, int(currentReservation.ID), orderDetails)
						if err == nil {
							orderId := int(orderDetails["id"].(float64))
							orderItems, err := r.GetRevelOrderItems(orderId, posIntegration)
							if err == nil {
								_, _ = r.CreateOrUpdateReservationOrderItemsFromRevel(ctx, reservationOrder.ID, orderItems)
							}
						}
					}
				}
			}
		}
	} else {
		updates["status_id"] = currentReservation.StatusID
	}

	if req.GetParams().GetDate() != "" &&
		currentReservation.Date.Format(time.DateOnly) != req.GetParams().GetDate() {
		reservationDate, err = time.Parse(time.DateOnly, req.GetParams().GetDate())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "date",
			OldValue:      currentReservation.Date.Format(time.DateOnly),
			NewValue:      req.GetParams().GetDate(),
		})
	} else {
		reservationDate = currentReservation.Date
	}

	if currentReservation.Date != reservationDate && req.GetParams().GetDate() != "" {
		updates["date"] = reservationDate
	}

	if req.GetParams().GetShiftId() != currentReservation.ShiftID && req.GetParams().GetShiftId() != 0 {
		updates["shift_id"] = req.GetParams().GetShiftId()

		if err := r.GetTenantDBConnection(ctx).
			Table("shifts").
			Where("id = ?", updates["shift_id"].(int32)).
			Select(`"id", "name", "from", "to", "time_interval"`).
			Scan(&shift).Error; err != nil {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrShiftNotFound)
		}

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "shift",
			OldValue:      currentReservation.Shift.Name,
			NewValue:      shift.Name,
		})
	} else {
		updates["shift_id"] = currentReservation.ShiftID
	}

	if req.GetParams().GetTime() != "" &&
		currentReservationTime.Format("15:04") != req.GetParams().GetTime() {
		if err := r.GetTenantDBConnection(ctx).
			Table("shifts").
			Where("id = ?", updates["shift_id"].(int32)).
			Select(`"id", "name", "from", "to", "time_interval"`).
			Scan(&shift).Error; err != nil {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrShiftNotFound)
		}

		reservationTime, err = time.Parse("15:04", req.GetParams().GetTime())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		fullTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), reservationTime.Hour(), reservationTime.Minute(), 0, 0, time.UTC)
		currentDay := time.Date(fullTime.Year(), fullTime.Month(), fullTime.Day(), 0, 0, 0, 0, time.UTC)
		dayBefore := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 0, 0, 0, 0, time.UTC)

		isTimeAfterMidnight := fullTime.Hour() >= 0 && fullTime.Hour() < 12 && shift.From.After(shift.To)

		if (!isTimeAfterMidnight || currentTime.Hour()*60+currentTime.Minute() > fullTime.Hour()*60+fullTime.Minute()) && currentDay.Equal(dayBefore) {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrInvalidTime)
		}

		if fullTime.Hour()*60+fullTime.Minute() > shift.To.Hour()*60+shift.To.Minute() &&
			fullTime.Hour()*60+fullTime.Minute() < shift.From.Hour()*60+shift.From.Minute() {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrInvalidTime)
		}

		if currentDay.Before(dayBefore) {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrInvalidDate)
		}

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "time",
			OldValue:      currentReservationTime.Format("15:04"),
			NewValue:      req.GetParams().GetTime(),
		})
	} else {
		reservationTime = currentReservationTime
	}

	if currentReservationTime != reservationTime && req.GetParams().GetTime() != "" {
		updates["time"] = reservationTime.Format("15:04")
	} else {
		updates["time"] = currentReservationTime.Format("15:04")
	}

	// Check if the new seating area is available
	if req.GetParams().GetSeatingAreaId() != 0 && req.GetParams().GetSeatingAreaId() != currentReservation.SeatingAreaID {
		updates["seating_area_id"] = req.GetParams().GetSeatingAreaId()

		seatingArea, err := r.seatingAreaClient.Client.SeatingAreaService.Repository.
			GetSeatingAreaByID(ctx, req.GetParams().GetSeatingAreaId())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrSeatingAreaNotFound)
		}

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "seating area",
			OldValue:      currentReservation.SeatingArea.Name,
			NewValue:      seatingArea.Name,
		})
	} else {
		updates["seating_area_id"] = currentReservation.SeatingAreaID
	}

	// Check if the new branch is available
	if req.GetParams().GetBranchId() != currentReservation.BranchID && req.GetParams().GetBranchId() != 0 {
		if !r.userClient.Client.UserService.Repository.CheckForBranchAccess(ctx, req.GetParams().GetBranchId()) {
			return nil, status.Error(http.StatusNotFound, "You don't have access to this branch")
		}

		updates["branch_id"] = req.GetParams().GetBranchId()

		var oldStatus string
		if err := r.GetTenantDBConnection(ctx).Table("reservation_statuses").Where("id = ?", updates["status_id"]).Select("name").Scan(&oldStatus).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		newStatus, err := r.GetReservationStatusByName(ctx, oldStatus, req.GetParams().GetBranchId())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		updates["status_id"] = int32(newStatus.ID)

		branch, err := r.userClient.Client.UserService.Repository.GetBranchByID(ctx, req.GetParams().GetBranchId())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "branch",
			OldValue:      currentReservation.Branch.Name,
			NewValue:      branch.Name,
		})
	} else {
		updates["branch_id"] = currentReservation.BranchID
	}

	if (req.GetParams().GetSeatingAreaId() != 0 && req.GetParams().GetSeatingAreaId() != currentReservation.SeatingAreaID) ||
		(req.GetParams().GetBranchId() != 0 && req.GetParams().GetBranchId() != currentReservation.BranchID) {
		tables, err := r.CreateReservationTables(ctx,
			reservationDate,
			reservationTime,
			reservationTime.Add(time.Minute*time.Duration(turnover.TurnoverTime)),
			int(updates["guests_number"].(int32)), int(updates["branch_id"].(int32)),
		)
		if err != nil && !isOverBooking {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
		}

		if err := r.GetTenantDBConnection(ctx).
			Delete(&domain2.ReservationTable{}, "reservation_id = ?", currentReservation.ID).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		reservationTable := &domain2.ReservationTable{}
		table := &tableDomain.Table{
			MaxPartySize: 1e9,
		}

		sameSeatingArea := false
		for _, t := range tables {
			if t.SeatingAreaID == int(updates["seating_area_id"].(int32)) {
				if t.MaxPartySize >= int(updates["guests_number"].(int32)) &&
					t.MinPartySize <= int(updates["guests_number"].(int32)) &&
					table.MaxPartySize >= t.MaxPartySize {
					table = t
				}
				sameSeatingArea = true
			}
		}

		if !sameSeatingArea && !isOverBooking {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrNoTablesAvailableInSeatingArea)
		}

		if table.MaxPartySize == 1e9 && !isOverBooking {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrMaxPartySize)
		}

		if sameSeatingArea && table.MaxPartySize != 1e9 {
			if table.CombinedTables != nil {
				for _, t := range commonUtils.ConvertStringToArrayBySeparator(*table.CombinedTables, ",") {
					tableId, _ := strconv.ParseInt(t, 10, 10)
					reservationTable = &domain2.ReservationTable{
						ReservationID: int(req.GetParams().GetId()),
						TableID:       int(tableId),
					}

					if err := r.GetTenantDBConnection(ctx).Create(reservationTable).Error; err != nil {
						return nil, status.Error(http.StatusInternalServerError, err.Error())
					}
				}
			} else {
				reservationTable = &domain2.ReservationTable{
					ReservationID: int(req.GetParams().GetId()),
					TableID:       table.ID,
				}

				if err := r.GetTenantDBConnection(ctx).Create(reservationTable).Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}
			}
		}

		updatedTime, err := time.Parse("15:04", updates["time"].(string))
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, "Invalid time format")
		}

		pacing, err := r.shiftClient.Client.ShiftService.Repository.GetPacingByTime(
			ctx,
			currentReservation.ShiftID,
			time.Date(reservationDate.Year(), reservationDate.Month(), reservationDate.Day(), updatedTime.Hour(), updatedTime.Minute(), 0, 0, time.UTC),
			currentReservation.BranchID,
			updates["seating_area_id"].(int32),
		)
		if err != nil {
			return nil, err
		}

		var oldPartySize int
		oldPartySize = int(currentReservation.GuestsNumber)
		if len(currentReservation.Tables) != 0 {
			oldPartySize = currentReservation.Tables[0].MaxPartySize
		}

		var capacity *shiftDomain.Pacing
		if err := r.GetTenantDBConnection(ctx).
			Model(&shiftDomain.Pacing{}).
			Where("shift_id = ? AND hour = ?",
				updates["shift_id"].(int32),
				updates["time"].(string)).
			Find(&capacity).
			Error; err != nil && !isOverBooking {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrMaxPartySize)
		}

		if capacity.Capacity < pacing+int32(oldPartySize-table.MaxPartySize) && !isOverBooking {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoTablesAvailable)
		}
	}

	currentStatus, err := r.CommonRepo.GetReservationStatusByID(ctx, updates["status_id"].(int32), updates["branch_id"].(int32))
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if req.GetParams().GetSeatedGuests() > updates["guests_number"].(int32) {
		return nil, status.Error(http.StatusNotFound, "Seated guests can't be greater than guests number")
	}

	if req.GetParams().GetSeatedGuests() != 0 &&
		currentStatus.Name == meta.PartiallySeated {
		currentReservation.SeatedGuests = req.GetParams().GetSeatedGuests()
		updates["seated_guests"] = req.GetParams().GetSeatedGuests()
	} else if currentStatus.Name == meta.Seated ||
		updates["guests_number"].(int32) == req.GetParams().GetGuestsNumber() {
		currentReservation.SeatedGuests = req.GetParams().GetGuestsNumber()
		updates["seated_guests"] = req.GetParams().GetGuestsNumber()
	}

	if req.GetParams().GetReservedVia() != currentReservation.ReservedVia && req.GetParams().GetReservedVia() != "" {
		updates["reserved_via"] = req.GetParams().GetReservedVia()

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "reserved via",
			OldValue:      currentReservation.ReservedVia,
			NewValue:      req.GetParams().GetReservedVia(),
		})
	}

	if req.GetParams().GetSpecialOccasionId() != 0 {
		updates["special_occasion_id"] = req.GetParams().GetSpecialOccasionId()

		specialOccasion, err := r.CommonRepo.GetReservationSpecialOccasionByID(ctx, req.GetParams().GetSpecialOccasionId())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		if currentReservation.SpecialOccasionID != nil {
			logs = append(logs, &reservationLogDomain.ReservationLog{
				ReservationID: req.GetParams().GetId(),
				FieldName:     "special occasion",
				OldValue:      currentReservation.SpecialOccasion.Name,
				NewValue:      specialOccasion.Name,
			})
		} else {
			logs = append(logs, &reservationLogDomain.ReservationLog{
				ReservationID: req.GetParams().GetId(),
				FieldName:     "special occasion",
				Action:        "create",
				NewValue:      specialOccasion.Name,
			})
		}
	}

	if req.GetParams().GetGuestId() != currentReservation.GuestID && req.GetParams().GetGuestId() != 0 {
		updates["guest_id"] = req.GetParams().GetGuestId()

		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "guest",
			OldValue:      strconv.FormatInt(int64(currentReservation.GuestID), 10),
			NewValue:      strconv.FormatInt(int64(req.GetParams().GetGuestId()), 10),
		})
	}

	if len(updates) > 0 {
		if err := r.GetTenantDBConnection(ctx).
			Model(&domain2.Reservation{}).
			Where("id = ?", req.GetParams().GetId()).
			Updates(updates).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	if !req.GetParams().GetEmptyTags() && utils.ArrayContains(permissions, rbac.EditReservationTags.Name) {
		var newTags []int32
		for _, tag := range req.GetParams().GetTags() {
			newTags = append(newTags, tag.Id)
		}

		var existingTags []int32
		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationTagsAssignment{}).
			Where("reservation_id = ?", req.GetParams().GetId()).
			Distinct("tag_id").
			Pluck("tag_id", &existingTags).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		oldTagsMap := make(map[int32]bool)
		for _, tag := range existingTags {
			oldTagsMap[tag] = true
		}

		newTagsMap := make(map[int32]bool)
		for _, tag := range newTags {
			newTagsMap[tag] = true
		}

		var deletedTags []string
		for tag := range oldTagsMap {
			if _, ok := newTagsMap[tag]; !ok {
				var fullTag *domain.ReservationTag

				if err := r.GetTenantDBConnection(ctx).
					Model(&domain.ReservationTag{}).
					Where("id =?", tag).
					Find(&fullTag).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				if err := r.GetTenantDBConnection(ctx).
					Model(&domain.ReservationTagsAssignment{}).
					Delete(tag, "tag_id = ?", tag).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				deletedTags = append(deletedTags, fullTag.Name)
			}
		}

		var createdTags []string
		for tag := range newTagsMap {
			if _, ok := oldTagsMap[tag]; !ok {
				if err := r.GetTenantDBConnection(ctx).
					Model(&domain.ReservationTagsAssignment{}).
					Create(&domain.ReservationTagsAssignment{
						ReservationID: req.GetParams().GetId(),
						TagID:         tag,
					}).Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				var fullTag *domain.ReservationTag
				if err := r.GetTenantDBConnection(ctx).
					Model(&domain.ReservationTag{}).
					Where("id =?", tag).
					First(&fullTag).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				createdTags = append(createdTags, fullTag.Name)
			}
		}

		if len(deletedTags) > 0 {
			var parsedDeletedTags string
			for i, tag := range deletedTags {
				if i == 0 {
					parsedDeletedTags = tag
					continue
				}
				parsedDeletedTags += ", " + tag
			}

			logs = append(logs, &reservationLogDomain.ReservationLog{
				ReservationID: req.GetParams().GetId(),
				FieldName:     "tags",
				Action:        "delete",
				OldValue:      parsedDeletedTags,
			})
		}

		if len(createdTags) > 0 {
			var parsedCreatedTags string
			for i, tag := range createdTags {
				if i == 0 {
					parsedCreatedTags = tag
					continue
				}
				parsedCreatedTags += ", " + tag
			}

			logs = append(logs, &reservationLogDomain.ReservationLog{
				ReservationID: req.GetParams().GetId(),
				FieldName:     "tags",
				Action:        "create",
				NewValue:      parsedCreatedTags,
			})
		}
	} else if !req.GetParams().GetEmptyTags() && !utils.ArrayContains(permissions, rbac.EditReservationTags.Name) {
		return nil, status.Error(http.StatusForbidden, "You don't have permission to edit reservation tags")
	}

	if req.GetParams().GetDeleteSpecialOccasion() {
		if err := r.GetTenantDBConnection(ctx).
			Model(&domain2.Reservation{}).
			Where("id = ?", req.GetParams().GetId()).
			Update("special_occasion_id", nil).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		logs = append(logs, &reservationLogDomain.ReservationLog{
			ReservationID: req.GetParams().GetId(),
			FieldName:     "special occasion",
			Action:        "delete",
		})
	}

	getUpdatedReservation, err := r.CommonRepo.GetReservationByID(ctx, req.GetParams().GetId())
	if err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrReservationNotFound)
	}

	for _, rlog := range logs {
		if rlog.Action == "" {
			rlog.Action = "update"
		}
	}
	if _, err := r.CreateReservationLogs(ctx, logs...); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	getUpdatedReservation.CheckOut = currentReservation.CheckOut

	if getUpdatedReservation.SpecialOccasionID != nil {
		if err := r.SendSpecialOccasionMessage(
			fmt.Sprintf("Special occasion %s has been added to reservation %d",
				getUpdatedReservation.SpecialOccasion.Name, getUpdatedReservation.ID),
			req.GetParams().GetId(),
		); err != nil {
			logger.Default().Error(err.Error())
		}
	}

	return &guestProto.UpdateReservationResponse{
		Result: converters.BuildReservationProto(getUpdatedReservation),
	}, nil
}
